package runner

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	runnerbuild "github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/debug"
	"github.com/xhd2015/doctest/libdoc/metrics"
	"github.com/xhd2015/doctest/libdoc/path_resolve"
	"github.com/xhd2015/doctest/libdoc/validate"
	"github.com/xhd2015/less-flags"
)

var ErrNoTestsFound = path_resolve.ErrNoTestsFound

// mergeRunStats accumulates per-tree stats into the CLI suite totals,
// including timeout Planned/TimedOut for cancelled FAIL presentation.
func mergeRunStats(dst *runnerbuild.TestRunStats, src runnerbuild.TestRunStats) {
	dst.Passed += src.Passed
	dst.Total += src.Total
	dst.SkipCount += src.SkipCount
	dst.Skipped = append(dst.Skipped, src.Skipped...)
	if src.Planned > 0 {
		dst.Planned += src.Planned
	}
	if src.TimedOut {
		dst.TimedOut = true
	}
	if src.BuildFailed {
		dst.BuildFailed = true
	}
	if src.GoTestBypassed {
		dst.GoTestBypassed = true
	}
	if src.NoTestsChanged {
		dst.NoTestsChanged = true
	}
	if dst.GenRoot == "" && src.GenRoot != "" {
		dst.GenRoot = src.GenRoot
	}
	if dst.TreeRel == "" && src.TreeRel != "" {
		dst.TreeRel = src.TreeRel
	}
}

// printGenPlanAfterRun emits plan + result trees from GenBatch after generate.
func printGenPlanAfterRun(w io.Writer, opts core.Options, remainArgs []string, multi bool, stats *runnerbuild.TestRunStats) {
	if w == nil {
		w = os.Stderr
	}
	genRoot := opts.GenDir
	if genRoot == "" && stats != nil && stats.GenRoot != "" {
		genRoot = stats.GenRoot
	}
	if genRoot == "" && opts.GenBatch != nil {
		if roots := opts.GenBatch.AllGenRoots(); len(roots) > 0 {
			genRoot = roots[0]
		}
	}
	planArgs := runnerbuild.ResolveGenPlanArgs(remainArgs, "")
	// Prefer TreeRel from stats for single-tree when Resolve fell back.
	if !multi && stats != nil && stats.TreeRel != "" && len(planArgs) == 1 {
		planArgs[0].TreeRel = filepath.ToSlash(stats.TreeRel)
	}
	runnerbuild.PrintGenPlanAndResult(w, opts, planArgs, genRoot, multi)
}

func Build(dir string) error {
	return runnerbuild.Build(dir, core.Options{RemoveTemp: false})
}

func BuildArgs(args []string) error {
	return BuildArgsWithWriters(args, nil, nil)
}

// BuildArgsWithWriters is BuildArgs with explicit stdout/stderr (no process globals).
func BuildArgsWithWriters(args []string, stdout, stderr io.Writer) error {
	return processArgsWithWriters(args, "build", parseBuildOptions, func(dir string, opts core.Options) error {
		err := runnerbuild.Build(dir, opts)
		if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
			return ErrNoTestsFound
		}
		return err
	}, stdout, stderr)
}

func Test(args []string) error {
	return TestWithWriters(args, nil, nil)
}

// TestWithWriters runs doctest test with optional stdout/stderr capture.
// Non-nil writers override parse defaults (including Stderr: os.Stderr).
// Safe for concurrent nested harnesses — no package-level inject state.
func TestWithWriters(args []string, stdout, stderr io.Writer) error {
	opts, remainArgs, err := parseTestOptions(args)
	if err != nil {
		return err
	}
	applyWriters(&opts, stdout, stderr)
	return runTest(opts, remainArgs)
}

// runTest is the body of Test after CLI flag parse. Same-package unit tests
// call it with opts.Stdout/Stderr injected (never swap process os.Stdout).
func runTest(opts core.Options, remainArgs []string) error {
	if len(remainArgs) < 1 {
		return fmt.Errorf("test requires <dir>")
	}

	// Resolve ColorAuto against the real user-facing stdout before any
	// parallel-tree buffering replaces opts.Stdout with a bytes.Buffer.
	// Otherwise Auto always disables color for `doctest test ./...`.
	stdoutForColor := opts.Stdout
	if stdoutForColor == nil {
		stdoutForColor = os.Stdout
	}
	opts.Color = runnerbuild.ResolveColorMode(opts.Color, stdoutForColor)

	// Engine-internal DOCTEST_DEBUG (GODEBUG-style). Fail closed on unknown keys.
	dbg, err := debug.FromEnv()
	if err != nil {
		return err
	}
	stderrW := opts.Stderr
	if stderrW == nil {
		stderrW = os.Stderr
	}
	if dbg.GenPlan {
		opts.GenPlan = true
		runnerbuild.PrintGenPlanBanner(stderrW)
	}
	if dbg.BypassGoTest {
		opts.BypassGoTest = true
		fmt.Fprintln(stderrW, "doctest: DOCTEST_DEBUG bypass-go-test=1 (go test will be skipped)")
	}

	// Cold-cache: resolve gen root, wipe on startup, force count, isolate GOCACHE.
	// Applied once per CLI invocation so multi-tree ./... shares GenDir/GOCACHE.
	var coldCacheNs int64
	if opts.ColdCache {
		tCold := time.Now()
		if err := applyColdCache(&opts); err != nil {
			return err
		}
		coldCacheNs = time.Since(tCold).Nanoseconds()
	}

	// -a: hard force — superset of -count=1 when count unset; gen wipe is
	// Options GenBatch.WipeOnce (not SessionID).
	applyForceA(&opts)

	// Shared gen batch for multi-tree emit-set union + -a wipe-once per gen root.
	if opts.GenBatch == nil {
		opts.GenBatch = core.NewGenBatch()
	}

	// gen-plan invocation header (after GenDir may be resolved by cold-cache).
	multiArg := len(remainArgs) > 1 || (len(remainArgs) == 1 && path_resolve.IsDotDotDotPattern(remainArgs[0]))
	if opts.GenPlan {
		genRootLabel := opts.GenDir
		runnerbuild.PrintGenPlanInvocation(stderrW, remainArgs, opts, genRootLabel, multiArg)
	}

	// One session id for the whole CLI invocation so parallel trees share
	// session.Once / testbin materialization when nested self-tests run.
	// Held on opts only — children receive it via cmd.Env key-replace (no process Setenv).
	// Unrelated to gen wipe / orphan prune.
	if opts.SessionID == "" {
		if v, ok := syscall.Getenv(core.DoctestSessionIDEnv); ok && v != "" {
			opts.SessionID = v
		} else {
			opts.SessionID = core.NewDoctestSessionID()
		}
	}

	// Honor DOCTEST_METRICS_ROOT when MetricsRoot not set via options.
	if opts.MetricsRoot == "" {
		if v := os.Getenv(EnvMetricsRoot); v != "" {
			opts.MetricsRoot = v
		}
	}

	opts.SuppressResultSummary = true
	start := time.Now()
	var stats runnerbuild.TestRunStats
	var statsMu sync.Mutex
	defaultSuite := !opts.LabelAll && len(opts.LabelExprs) == 0

	// Suite-level metrics for the whole CLI invocation (one JSONL when a
	// single tree; still one file for multi-arg — first tree path as root).
	// Opt-in only: --metrics-on.
	var rec *runRecorder
	if opts.MetricsOn {
		metricDir := remainArgs[0]
		if path_resolve.IsDotDotDotPattern(metricDir) {
			metricDir = path_resolve.ExtractBasePath(metricDir)
		}
		if metricDir == "" || metricDir == "..." {
			metricDir, _ = os.Getwd()
		}
		// resolve target for project id
		if !path_resolve.IsDotDotDotPattern(remainArgs[0]) {
			td, _ := resolveTestTarget(remainArgs[0])
			if root, ok := path_resolve.ResolveRoot(td); ok {
				metricDir = root
			} else {
				metricDir = td
			}
		}
		var openErr error
		rec, openErr = openRunRecorder(metricDir, opts)
		if openErr != nil {
			fmt.Fprintf(opts.Stderr, "doctest: metrics: %v\n", openErr)
			rec = nil
		}
		if rec != nil {
			defer rec.close()
			_ = rec.writeRunStart(metricDir, opts, defaultSuite)
			if coldCacheNs > 0 {
				_ = rec.writePhase("suite", "cold_cache", "", coldCacheNs, nil)
			}
			// Nest sink for go test children via Options → goTestEnv (cmd.Env only).
			nestSink := rec.path + ".nest"
			_ = os.Remove(nestSink)
			opts.MetricsNestSink = nestSink
			rec.nestSink = nestSink
			defer func() {
				_ = os.Remove(nestSink)
			}()
		}
	}

	runFn := func(dir string, o core.Options) error {
		o.SuppressResultSummary = true
		// Avoid nested per-tree files; CLI owns the suite recorder above.
		o.MetricsOn = false
		// Leaf events when recorder is active
		var cases []core.TreeCase
		if rec != nil {
			if cs, err := core.DiscoverTreeCases(dir); err == nil {
				if o.SubDir != "" {
					cs = core.FilterBySubDir(cs, dir, o.SubDir)
				}
				cs, _ = core.FilterCasesByLabel(cs, o)
				cases = cs
				for _, c := range cases {
					_ = rec.writeLeafStart(c, dir)
				}
				s, err := runnerbuild.TestWithStats(dir, o)
				// Tree phase spans (discover/materialize/generate/go_test/…).
				for _, p := range s.Phases {
					_ = rec.writePhase("tree", p.Name, dir, p.ElapsedNs, map[string]any{
						"cases": s.Total,
					})
				}
				// Leaf ends: prefer package-attributed / unified subtest times.
				timingByPath := map[string]runnerbuild.LeafTiming{}
				for _, lt := range s.LeafTimings {
					timingByPath[lt.Path] = lt
				}
				passLeft := s.Passed
				end := time.Now()
				for _, c := range cases {
					result := "fail"
					if s.GoTestBypassed {
						result = "bypassed"
					} else if passLeft > 0 {
						result = "pass"
						passLeft--
					}
					lt, ok := timingByPath[c.Path]
					var elapsedNs int64
					cached := false
					if ok {
						elapsedNs = lt.ElapsedNs
						cached = lt.Cached
					}
					_ = rec.writeLeafEndNs(c, dir, end, elapsedNs, result, cached)
				}
				for _, sk := range s.Skipped {
					_ = rec.writeLeafEndSkipped(sk, dir, end)
				}
				statsMu.Lock()
				mergeRunStats(&stats, s)
				statsMu.Unlock()
				if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
					return ErrNoTestsFound
				}
				return err
			}
		}
		s, err := runnerbuild.TestWithStats(dir, o)
		statsMu.Lock()
		mergeRunStats(&stats, s)
		statsMu.Unlock()
		if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
			return ErrNoTestsFound
		}
		return err
	}

	var runErr error
	if len(remainArgs) == 1 && path_resolve.IsDotDotDotPattern(remainArgs[0]) {
		// Multi-root: generate all unified trees, then one go test on
		// __workspace/suite (workspace registry + per-tree __wreg).
		runErr = testDotDotDotWorkspace(remainArgs[0], opts, rec, &stats, &statsMu)
	} else if len(remainArgs) == 1 {
		// Single-tree: TestWithStats (prepareLeafCache + go test); same leaf-cache policy.
		runErr = processSingleArg(remainArgs[0], opts, runFn)
	} else {
		// P3: multi-arg shares DiscoverRoots → PrepareAll → BindExec → RunSuite
		// with ./... when roots can bind one workspace/hub. Conflicting filters
		// on the same root fall back to N× TestWithStats.
		targets, expErr := expandTestArgs(remainArgs)
		switch {
		case expErr != nil:
			runErr = expErr
		case suiteTargetsConflict(targets):
			runErr = processMultiArg(remainArgs, opts, runFn)
		default:
			runErr = runSuitePlan(targets, opts, rec, &stats, &statsMu)
		}
	}

	// gen-plan plan + result trees after emit/prune are known.
	// Printed here (not inside PrepareTree) so multi-arg buffered prepare
	// stderr is not swallowed when -v is off.
	if opts.GenPlan {
		printGenPlanAfterRun(stderrW, opts, remainArgs, multiArg, &stats)
	}

	if len(stats.Skipped) > 0 {
		sw := opts.Stdout
		if sw == nil {
			sw = os.Stdout
		}
		runnerbuild.PrintSkippedSummaryTo(sw, stats.Skipped, opts.Verbose)
	}
	if stats.Total > 0 || stats.SkipCount > 0 || stats.BuildFailed {
		stats.Elapsed = time.Since(start)
		// Soft "no tests" is not an overall failure for the summary line.
		// Any other runErr (prepare, go test, multi-tree) must not print PASS
		// even when survivor cases all passed.
		overallOK := runErr == nil || errors.Is(runErr, ErrNoTestsFound)
		runnerbuild.PrintResultSummaryOverall(opts, stats, overallOK)
	}

	// Close metrics with run_end after summary.
	if rec != nil {
		_ = mergeNestSinkIntoRecorder(rec, rec.nestSink)
		warns := []string{}
		if metrics.ShouldWarnDefaultSuiteSlow(defaultSuite, stats.Total, stats.Elapsed, metrics.DefaultSuiteWarnThreshold) {
			warns = append(warns, "default_suite_slow")
		}
		exitOK := runErr == nil && stats.Total > 0 && (stats.GoTestBypassed || stats.Passed >= stats.Total)
		_ = rec.writeRunEnd(stats, exitOK, warns)
		_ = rec.close()
		rec = nil
	}

	if metrics.ShouldWarnDefaultSuiteSlow(defaultSuite, stats.Total, stats.Elapsed, metrics.DefaultSuiteWarnThreshold) {
		fmt.Fprint(opts.Stderr, metrics.FormatDefaultSuiteSlowWarning())
	}

	if runErr != nil {
		if errors.Is(runErr, ErrNoTestsFound) && stats.Total == 0 && len(stats.Skipped) == 0 {
			return ErrNoTestsFound
		}
		if errors.Is(runErr, ErrNoTestsFound) && len(stats.Skipped) > 0 {
			return nil
		}
		return runErr
	}
	if stats.Total == 0 {
		// Runtime t.Skip alone (actual_run=0) is not a suite failure.
		if stats.NoTestsChanged || len(stats.Skipped) > 0 || stats.SkipCount > 0 {
			return nil
		}
		return ErrNoTestsFound
	}
	if stats.GoTestBypassed {
		return nil
	}
	if stats.Passed < stats.Total {
		return fmt.Errorf("%d of %d tests passed", stats.Passed, stats.Total)
	}
	return nil
}

func VetArgs(args []string) error {
	return VetArgsWithWriters(args, nil, nil)
}

// VetArgsWithWriters is VetArgs with explicit stdout/stderr (no process globals).
func VetArgsWithWriters(args []string, stdout, stderr io.Writer) error {
	return processArgsWithWriters(args, "vet", parseVetOptions, func(dir string, opts core.Options) error {
		return validate.RunWithOptions(dir, opts)
	}, stdout, stderr)
}

func processArgs(args []string, cmdName string, parseFn func([]string) (core.Options, []string, error), processDirFn func(string, core.Options) error) error {
	return processArgsWithWriters(args, cmdName, parseFn, processDirFn, nil, nil)
}

func processArgsWithWriters(args []string, cmdName string, parseFn func([]string) (core.Options, []string, error), processDirFn func(string, core.Options) error, stdout, stderr io.Writer) error {
	opts, remainArgs, err := parseFn(args)
	if err != nil {
		return err
	}
	applyWriters(&opts, stdout, stderr)
	if len(remainArgs) < 1 {
		return fmt.Errorf("%s requires <dir>", cmdName)
	}
	if len(remainArgs) == 1 {
		return processSingleArg(remainArgs[0], opts, processDirFn)
	}
	return processMultiArg(remainArgs, opts, processDirFn)
}

// suiteTarget is one DOCTEST root (optional SubDir / ExplicitLeaf filter) in a
// suite plan: DiscoverRoots → PrepareAll → BindExec → RunSuite → Report.
type suiteTarget struct {
	Root         string
	SubDir       string
	ExplicitLeaf bool
}

// expandTestArgs expands CLI remain args into a deduped union of suite targets.
// Supports explicit paths and ./... patterns (same discovery as ./... workspace).
func expandTestArgs(args []string) ([]suiteTarget, error) {
	var targets []suiteTarget
	seen := map[string]bool{}
	add := func(t suiteTarget) {
		r := filepath.Clean(t.Root)
		s := filepath.Clean(t.SubDir)
		key := r + "\x00" + s + "\x00"
		if t.ExplicitLeaf {
			key += "1"
		} else {
			key += "0"
		}
		if seen[key] {
			return
		}
		seen[key] = true
		t.Root = r
		t.SubDir = s
		targets = append(targets, t)
	}

	for _, arg := range args {
		if arg == "..." {
			return nil, fmt.Errorf("bare '...' pattern is not supported; use './...' or 'path/...' instead")
		}
		if path_resolve.IsDotDotDotPattern(arg) {
			dirs, err := path_resolve.FindDotDotDotDirs(path_resolve.ExtractBasePath(arg))
			if err != nil {
				return nil, err
			}
			for _, d := range dirs {
				root, ok := path_resolve.ResolveRoot(d)
				if !ok || root == "" {
					root = d
				}
				// SubDir is the discovered path d, not always root.
				// Mid-tree path/... (e.g. suite-selection/...) must filter parent
				// DOCTEST leaves to under d; when d is itself a DOCTEST root,
				// ResolveRoot(d)==d so SubDir==root (full tree for that root).
				add(suiteTarget{Root: root, SubDir: d, ExplicitLeaf: false})
			}
			continue
		}
		targetDir, explicitLeaf := resolveTestTarget(arg)
		root, ok := path_resolve.ResolveRoot(targetDir)
		if !ok {
			root = targetDir
		}
		add(suiteTarget{Root: root, SubDir: targetDir, ExplicitLeaf: explicitLeaf})
	}
	return targets, nil
}

// suiteTargetsConflict is true when the same Root appears with different
// SubDir/ExplicitLeaf filters — cannot safely bind one workspace suite.
func suiteTargetsConflict(targets []suiteTarget) bool {
	type filt struct {
		subDir string
		leaf   bool
		set    bool
	}
	byRoot := map[string]filt{}
	for _, t := range targets {
		r := filepath.Clean(t.Root)
		f := filt{subDir: filepath.Clean(t.SubDir), leaf: t.ExplicitLeaf, set: true}
		if prev, ok := byRoot[r]; ok && prev.set {
			if prev.subDir != f.subDir || prev.leaf != f.leaf {
				return true
			}
			continue
		}
		byRoot[r] = f
	}
	return false
}

// testDotDotDotWorkspace discovers DOCTEST roots under a ./... pattern and
// runs the shared suite plan (PrepareAll + RunWorkspace / legacy per-root).
func testDotDotDotWorkspace(arg string, opts core.Options, rec *runRecorder, stats *runnerbuild.TestRunStats, statsMu *sync.Mutex) error {
	targets, err := expandTestArgs([]string{arg})
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return ErrNoTestsFound
	}
	return runSuitePlan(targets, opts, rec, stats, statsMu)
}

// runSuitePlan is the shared CLI plan path for multi-root invocations
// (./... and multi-arg when non-conflicting):
//
//	PrepareAll (parallel PrepareTree) → BindExec (RunWorkspace for unified;
//	TestWithStats for non-unified) → aggregate stats for suite Report.
//
// Metrics recorder stays suite-level (one file); nested trees set MetricsOn=false.
func runSuitePlan(targets []suiteTarget, opts core.Options, rec *runRecorder, stats *runnerbuild.TestRunStats, statsMu *sync.Mutex) error {
	if len(targets) == 0 {
		return ErrNoTestsFound
	}

	// Parallel prepare (generate-only), same worker pool as RunForDirs.
	type prepResult struct {
		prep         runnerbuild.TreePrep
		err          error
		dir          string
		subDir       string
		explicitLeaf bool
	}
	results := make([]prepResult, len(targets))
	var wg sync.WaitGroup
	// Bound concurrency similar to defaultTreeWorkers.
	workers := runtime.NumCPU()
	if workers < 4 {
		workers = 4
	}
	if workers > 12 {
		workers = 12
	}
	sem := make(chan struct{}, workers)
	for i, tgt := range targets {
		wg.Add(1)
		go func(i int, tgt suiteTarget) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			o := opts
			o.SubDir = tgt.SubDir
			o.ExplicitLeaf = tgt.ExplicitLeaf
			o.MetricsOn = false
			o.SuppressResultSummary = true
			// Isolate noisy per-tree headers into buffers; workspace run prints once.
			var outBuf, errBuf bytes.Buffer
			o.Stdout = &outBuf
			o.Stderr = &errBuf
			prep, err := runnerbuild.PrepareTree(tgt.Root, o)
			results[i] = prepResult{prep: prep, err: err, dir: tgt.Root, subDir: tgt.SubDir, explicitLeaf: tgt.ExplicitLeaf}
			// Drop generate chatter unless verbose (user sees workspace go test line).
			if opts.Verbose {
				stderrBase := opts.Stderr
				if stderrBase == nil {
					stderrBase = os.Stderr
				}
				stdoutBase := opts.Stdout
				if stdoutBase == nil {
					stdoutBase = os.Stdout
				}
				_, _ = io.Copy(stderrBase, &errBuf)
				_, _ = io.Copy(stdoutBase, &outBuf)
			}
		}(i, tgt)
	}
	wg.Wait()

	var (
		unified  []runnerbuild.TreePrep
		legacy   []prepResult
		prepErrs []string
		runErrs  []string
	)
	for _, r := range results {
		if r.err != nil {
			if errors.Is(r.err, ErrNoTestsFound) || strings.Contains(r.err.Error(), "no runnable test cases found") {
				// Skipped-only trees still contribute skipped stats.
				if len(r.prep.Skipped) > 0 {
					statsMu.Lock()
					stats.Skipped = append(stats.Skipped, r.prep.Skipped...)
					statsMu.Unlock()
				}
				continue
			}
			if r.prep.Stats.NoTestsChanged {
				statsMu.Lock()
				stats.NoTestsChanged = true
				statsMu.Unlock()
				continue
			}
			prepErrs = append(prepErrs, r.dir+": "+r.err.Error())
			continue
		}
		if r.prep.Stats.NoTestsChanged {
			statsMu.Lock()
			stats.NoTestsChanged = true
			statsMu.Unlock()
			continue
		}
		if r.prep.Stats.Total == 0 {
			if len(r.prep.Skipped) > 0 {
				statsMu.Lock()
				stats.Skipped = append(stats.Skipped, r.prep.Skipped...)
				statsMu.Unlock()
			}
			continue
		}
		if r.prep.Unified {
			unified = append(unified, r.prep)
		} else {
			legacy = append(legacy, r)
		}
	}

	// Metrics: leaf_start for all unified cases before the single go test.
	if rec != nil {
		for _, p := range unified {
			for _, c := range p.Cases {
				_ = rec.writeLeafStart(c, p.AbsRoot)
			}
			for _, ph := range p.Stats.Phases {
				_ = rec.writePhase("tree", ph.Name, p.AbsRoot, ph.ElapsedNs, map[string]any{
					"cases": p.Stats.Total,
				})
			}
		}
	}

	if len(unified) > 0 {
		wsOpts := opts
		wsOpts.MetricsOn = false
		wsOpts.SuppressResultSummary = true
		s, wsErr := runnerbuild.RunWorkspace(unified, wsOpts)
		statsMu.Lock()
		mergeRunStats(stats, s)
		statsMu.Unlock()
		if rec != nil {
			for _, ph := range s.Phases {
				if ph.Name == "go_test" {
					detail := map[string]any{
						"trees": len(unified),
						"cases": s.Total,
					}
					if s.GoTestBypassed {
						detail["bypassed"] = true
					}
					_ = rec.writePhase("suite", "go_test", "", ph.ElapsedNs, detail)
				}
			}
			// Attribute leaf ends; timings may be sparse when paths collide across trees.
			timingByPath := map[string]runnerbuild.LeafTiming{}
			for _, lt := range s.LeafTimings {
				timingByPath[lt.Path] = lt
			}
			end := time.Now()
			// Parallel trees make per-leaf pass/fail ordering unreliable; mark
			// pass when the workspace go test succeeded.
			for _, p := range unified {
				for _, c := range p.Cases {
					result := "fail"
					if s.GoTestBypassed {
						result = "bypassed"
					} else if wsErr == nil {
						result = "pass"
					}
					lt := timingByPath[c.Path]
					_ = rec.writeLeafEndNs(c, p.AbsRoot, end, lt.ElapsedNs, result, lt.Cached)
				}
				for _, sk := range p.Skipped {
					_ = rec.writeLeafEndSkipped(sk, p.AbsRoot, end)
				}
			}
		}
		if wsErr != nil {
			runErrs = append(runErrs, wsErr.Error())
		}
	}

	// Non-unified (internal-compile) trees: legacy per-root go test.
	for _, r := range legacy {
		o := opts
		o.SubDir = r.subDir
		o.ExplicitLeaf = r.explicitLeaf
		o.MetricsOn = false
		o.SuppressResultSummary = true
		s, err := runnerbuild.TestWithStats(r.dir, o)
		statsMu.Lock()
		mergeRunStats(stats, s)
		statsMu.Unlock()
		if err != nil && !strings.Contains(err.Error(), "no runnable test cases found") {
			runErrs = append(runErrs, r.dir+": "+err.Error())
		}
	}

	if err := formatClassifiedErrors(prepErrs, runErrs); err != nil {
		return err
	}
	if stats.Total == 0 && len(stats.Skipped) == 0 && !stats.NoTestsChanged {
		return ErrNoTestsFound
	}
	return nil
}

// formatClassifiedErrors builds multi-tree failure text with honest labels:
// prepare-only → "prepare failed:", go-test/workspace-only → "test failures:",
// both → "errors:".
func formatClassifiedErrors(prepErrs, runErrs []string) error {
	if len(prepErrs) == 0 && len(runErrs) == 0 {
		return nil
	}
	sort.Strings(prepErrs)
	sort.Strings(runErrs)
	switch {
	case len(prepErrs) > 0 && len(runErrs) > 0:
		all := append(append([]string(nil), prepErrs...), runErrs...)
		return fmt.Errorf("errors:\n%s", strings.Join(all, "\n"))
	case len(prepErrs) > 0:
		return fmt.Errorf("prepare failed:\n%s", strings.Join(prepErrs, "\n"))
	default:
		return fmt.Errorf("test failures:\n%s", strings.Join(runErrs, "\n"))
	}
}

func processSingleArg(arg string, opts core.Options, fn func(string, core.Options) error) error {
	if arg == "..." {
		return fmt.Errorf("bare '...' pattern is not supported; use './...' or 'path/...' instead")
	}
	if path_resolve.IsDotDotDotPattern(arg) {
		// Fallback multi-tree path (build/vet, or if testDotDotDot not used).
		// Parallel trees: buffer each tree's streams so progress lines do not interleave.
		var printMu sync.Mutex
		stdoutBase := opts.Stdout
		if stdoutBase == nil {
			stdoutBase = os.Stdout
		}
		stderrBase := opts.Stderr
		if stderrBase == nil {
			stderrBase = os.Stderr
		}
		return path_resolve.RunForDirs(path_resolve.ExtractBasePath(arg), func(dir string) error {
			root, _ := path_resolve.ResolveRoot(dir)
			if root == "" {
				root = dir
			}
			var outBuf, errBuf bytes.Buffer
			o := opts
			o.SubDir = dir
			o.ExplicitLeaf = false
			o.Stdout = &outBuf
			o.Stderr = &errBuf
			err := fn(root, o)
			// stderr first: tree header / "cd ..." then progress on stdout
			printMu.Lock()
			_, _ = io.Copy(stderrBase, &errBuf)
			_, _ = io.Copy(stdoutBase, &outBuf)
			printMu.Unlock()
			return err
		})
	}
	targetDir, explicitLeaf := resolveTestTarget(arg)
	root, ok := path_resolve.ResolveRoot(targetDir)
	if !ok {
		root = targetDir
	}
	opts.SubDir = targetDir
	opts.ExplicitLeaf = explicitLeaf
	return fn(root, opts)
}

func processMultiArg(args []string, opts core.Options, fn func(string, core.Options) error) error {
	var errs []string
	allNoTestsFound := true

	for _, arg := range args {
		if arg == "..." {
			return fmt.Errorf("bare '...' pattern is not supported; use './...' or 'path/...' instead")
		}
		if path_resolve.IsDotDotDotPattern(arg) {
			err := path_resolve.RunForDirs(path_resolve.ExtractBasePath(arg), func(dir string) error {
				root, _ := path_resolve.ResolveRoot(dir)
				if root == "" {
					root = dir
				}
				o := opts
				o.SubDir = dir
				o.ExplicitLeaf = false
				return fn(root, o)
			})
			if err != nil {
				if errors.Is(err, ErrNoTestsFound) {
					continue
				}
				errs = append(errs, err.Error())
				allNoTestsFound = false
			} else {
				allNoTestsFound = false
			}
			continue
		}
		targetDir, explicitLeaf := resolveTestTarget(arg)
		root, ok := path_resolve.ResolveRoot(targetDir)
		if !ok {
			root = targetDir
		}
		o := opts
		o.SubDir = targetDir
		o.ExplicitLeaf = explicitLeaf
		err := fn(root, o)
		if errors.Is(err, ErrNoTestsFound) {
			continue
		}
		if err != nil {
			errs = append(errs, err.Error())
			allNoTestsFound = false
		} else {
			allNoTestsFound = false
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("test failures:\n%s", strings.Join(errs, "\n"))
	}
	if allNoTestsFound {
		return ErrNoTestsFound
	}
	return nil
}

func parseBuildOptions(args []string) (core.Options, []string, error) {
	opts := core.Options{Stderr: os.Stderr, RemoveTemp: false}
	remainArgs, err := lessflags.Bool("-v,--verbose", &opts.Verbose).
		Bool("--rm", &opts.RemoveTemp).
		String("--gen-dir", &opts.GenDir).
		Int("-count", &opts.Count).
		Bool("--changed", &opts.ChangedOnly).
		Parse(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	return opts, remainArgs, nil
}

func extractLabelFlags(args []string) (labelExprs []string, remain []string, err error) {
	for i := 0; i < len(args); i++ {
		if args[i] == "--label" {
			if i+1 >= len(args) {
				return nil, nil, fmt.Errorf("--label requires an expression argument")
			}
			labelExprs = append(labelExprs, args[i+1])
			i++
			continue
		}
		remain = append(remain, args[i])
	}
	return labelExprs, remain, nil
}

func parseTestOptions(args []string) (core.Options, []string, error) {
	labelExprs, args, err := extractLabelFlags(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	for _, expr := range labelExprs {
		if err := core.ParseLabelExpr(expr); err != nil {
			return core.Options{}, nil, err
		}
	}

	var sawColor, sawNoColor bool
	for _, arg := range args {
		if arg == "--color" {
			sawColor = true
		}
		if arg == "--no-color" {
			sawNoColor = true
		}
	}
	if sawColor && sawNoColor {
		return core.Options{}, nil, fmt.Errorf("--color and --no-color are mutually exclusive")
	}

	opts := core.Options{Stderr: os.Stderr, RemoveTemp: false, Color: core.ColorAuto}
	var colorFlag, noColorFlag bool
	// **time.Duration: nil = -timeout omitted (go default 10m); non-nil 0 = disable.
	// Primary form is -timeout (go test style); --timeout kept as alias.
	var timeout *time.Duration
	remainArgs, err := lessflags.Bool("-v,--verbose", &opts.Verbose).
		Bool("--rm", &opts.RemoveTemp).
		String("--gen-dir", &opts.GenDir).
		Int("-count", &opts.Count).
		Bool("-a", &opts.ForceWithFlagA).
		Duration("-timeout,--timeout", &timeout).
		Bool("--color", &colorFlag).
		Bool("--no-color", &noColorFlag).
		Bool("--changed", &opts.ChangedOnly).
		Bool("--label-all", &opts.LabelAll).
		Bool("--metrics-on", &opts.MetricsOn).
		Bool("--cold-cache", &opts.ColdCache).
		String("-cpuprofile", &opts.CPUProfile).
		String("-memprofile", &opts.MemProfile).
		Int("-memprofilerate", &opts.MemProfileRate).
		String("-blockprofile", &opts.BlockProfile).
		Int("-blockprofilerate", &opts.BlockProfileRate).
		String("-mutexprofile", &opts.MutexProfile).
		Int("-mutexprofilefraction", &opts.MutexProfileFraction).
		String("-trace", &opts.Trace).
		String("-outputdir", &opts.OutputDir).
		String("-coverprofile", &opts.CoverProfile).
		Bool("-cover", &opts.Cover).
		Bool("-race", &opts.Race).
		String("--go-cmd", &opts.GoCmd).
		Parse(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	if err := core.ValidateGoCmdMode(opts.GoCmd); err != nil {
		return core.Options{}, nil, err
	}
	if colorFlag {
		opts.Color = core.ColorAlways
	}
	if noColorFlag {
		opts.Color = core.ColorNever
	}
	opts.Timeout = timeout
	if opts.LabelAll && len(labelExprs) > 0 {
		return core.Options{}, nil, fmt.Errorf("--label-all and --label are mutually exclusive")
	}
	opts.LabelExprs = labelExprs

	// Abs-resolve relative profile/cover paths against process cwd at parse time.
	pathFields := []*string{
		&opts.CPUProfile,
		&opts.MemProfile,
		&opts.BlockProfile,
		&opts.MutexProfile,
		&opts.Trace,
		&opts.OutputDir,
		&opts.CoverProfile,
	}
	for _, p := range pathFields {
		if *p == "" {
			continue
		}
		if !filepath.IsAbs(*p) {
			abs, absErr := absProfilePath(*p)
			if absErr != nil {
				return core.Options{}, nil, absErr
			}
			*p = abs
		}
	}

	// Ensure parent directories exist so go test can create profile files.
	// OutputDir itself is a directory destination.
	for _, dirPath := range []string{opts.OutputDir} {
		if dirPath == "" {
			continue
		}
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return core.Options{}, nil, fmt.Errorf("create -outputdir: %w", err)
		}
	}
	for _, filePath := range []string{
		opts.CPUProfile,
		opts.MemProfile,
		opts.BlockProfile,
		opts.MutexProfile,
		opts.Trace,
		opts.CoverProfile,
	} {
		if filePath == "" {
			continue
		}
		if dir := filepath.Dir(filePath); dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return core.Options{}, nil, fmt.Errorf("create profile parent dir: %w", err)
			}
		}
	}

	return opts, remainArgs, nil
}

// absProfilePath resolves a relative path against the process cwd.
// On macOS, Getwd often returns the /private/var form while user-facing
// TempDir paths use /var/...; prefer the non-/private form for stable matching.
func absProfilePath(p string) (string, error) {
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	const priv = "/private"
	if strings.HasPrefix(abs, priv+"/var/") {
		return abs[len(priv):], nil
	}
	return abs, nil
}

// ParseTestOptions parses doctest test flags into core.Options.
// Exported for package-level tests (metrics flags, labels, etc.).
func ParseTestOptions(args []string) (core.Options, []string, error) {
	return parseTestOptions(args)
}

// applyForceA resolves -a once per CLI invocation: when count is unset, force
// -count=1 (superset of count-based leaf-cache / go testcache busting). Gen wipe
// and go test -a are applied at generate / go test time via ForceWithFlagA.
func applyForceA(opts *core.Options) {
	if opts == nil || !opts.ForceWithFlagA {
		return
	}
	if opts.Count == 0 {
		opts.Count = 1
	}
}

// applyColdCache resolves --cold-cache semantics once per CLI invocation:
// force count=1 when unset, choose/protect gen root, wipe on startup only,
// isolate GOCACHE, and announce on stderr.
func applyColdCache(opts *core.Options) error {
	if opts == nil || !opts.ColdCache {
		return nil
	}
	if opts.Count == 0 {
		opts.Count = 1
	}

	cacheHome, err := core.CacheHome()
	if err != nil {
		return fmt.Errorf("cold-cache: resolve cache home: %w", err)
	}
	warmHome := filepath.Clean(filepath.Join(cacheHome, "doctest", "mapping-gen"))
	coldHome := filepath.Clean(filepath.Join(cacheHome, "doctest", "mapping-gen-cold"))
	if abs, absErr := filepath.Abs(warmHome); absErr == nil {
		warmHome = filepath.Clean(abs)
	}
	if abs, absErr := filepath.Abs(coldHome); absErr == nil {
		coldHome = filepath.Clean(abs)
	}

	gen := opts.GenDir
	if gen == "" {
		gen = coldHome
	} else {
		absGen, absErr := filepath.Abs(gen)
		if absErr != nil {
			return fmt.Errorf("cold-cache: resolve --gen-dir: %w", absErr)
		}
		gen = filepath.Clean(absGen)
		if pathEqualOrUnder(gen, warmHome) {
			// Do not wipe warm mapping-gen content on reject.
			return fmt.Errorf("error: --cold-cache refuses --gen-dir equal to or under warm mapping-gen (%s); cannot use the default warm cache path — choose mapping-gen-cold or another path outside %s", gen, warmHome)
		}
	}

	if err := os.RemoveAll(gen); err != nil {
		return fmt.Errorf("cold-cache: wipe gen dir %s: %w", gen, err)
	}
	core.InvalidateGenManifestCache(gen)
	if err := os.MkdirAll(gen, 0o755); err != nil {
		return fmt.Errorf("cold-cache: create gen dir %s: %w", gen, err)
	}
	opts.GenDir = gen

	gocacheTemp, err := os.MkdirTemp("", "doctest-cold-gocache-*")
	if err != nil {
		return fmt.Errorf("cold-cache: create isolated GOCACHE: %w", err)
	}
	opts.GoCache = gocacheTemp
	// Isolation is via opts.GoCache → child go tool cmd.Env (key-replace).
	// Do not mutate process env for GOCACHE (races under parallel suite leaves).

	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, "doctest: cold-cache: gen=%s GOCACHE=%s (isolated) count=%d\n", gen, gocacheTemp, opts.Count)
	return nil
}

// pathEqualOrUnder reports whether path is equal to root or a path under root.
// Both paths should already be cleaned absolute paths when possible.
func pathEqualOrUnder(path, root string) bool {
	path = filepath.Clean(path)
	root = filepath.Clean(root)
	if path == root {
		return true
	}
	sep := string(filepath.Separator)
	return strings.HasPrefix(path, root+sep)
}

func parseVetOptions(args []string) (core.Options, []string, error) {
	opts := core.Options{Stderr: os.Stderr}
	remainArgs, err := lessflags.Bool("-v,--verbose", &opts.Verbose).
		Bool("--changed", &opts.ChangedOnly).
		Parse(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	return opts, remainArgs, nil
}
