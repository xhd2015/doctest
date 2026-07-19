package runner

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/xhd2015/less-flags"
	runnerbuild "github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/metrics"
	"github.com/xhd2015/doctest/libdoc/path_resolve"
	"github.com/xhd2015/doctest/libdoc/validate"
)

var ErrNoTestsFound = path_resolve.ErrNoTestsFound

func Build(dir string) error {
	return runnerbuild.Build(dir, core.Options{RemoveTemp: false})
}

func BuildArgs(args []string) error {
	return processArgs(args, "build", parseBuildOptions, func(dir string, opts core.Options) error {
		err := runnerbuild.Build(dir, opts)
		if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
			return ErrNoTestsFound
		}
		return err
	})
}

func Test(args []string) error {
	opts, remainArgs, err := parseTestOptions(args)
	if err != nil {
		return err
	}
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

	// One session id for the whole CLI invocation so parallel trees share
	// session.Once / testbin materialization when nested self-tests run.
	if v, ok := syscall.Getenv(core.DoctestSessionIDEnv); !ok || v == "" {
		_ = os.Setenv(core.DoctestSessionIDEnv, core.NewDoctestSessionID())
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
			// Nest sink: suite children (nested RunTest) append timing here;
			// merged into the outer JSONL before run_end (see below).
			nestSink := rec.path + ".nest"
			_ = os.Remove(nestSink)
			_ = os.Setenv(metrics.EnvMetricsNestSink, nestSink)
			// Stash on recorder for merge before close.
			rec.nestSink = nestSink
			defer func() {
				_ = os.Unsetenv(metrics.EnvMetricsNestSink)
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
					if passLeft > 0 {
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
				stats.Passed += s.Passed
				stats.Total += s.Total
				stats.Skipped = append(stats.Skipped, s.Skipped...)
				if s.NoTestsChanged {
					stats.NoTestsChanged = true
				}
				statsMu.Unlock()
				if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
					return ErrNoTestsFound
				}
				return err
			}
		}
		s, err := runnerbuild.TestWithStats(dir, o)
		statsMu.Lock()
		stats.Passed += s.Passed
		stats.Total += s.Total
		stats.Skipped = append(stats.Skipped, s.Skipped...)
		if s.NoTestsChanged {
			stats.NoTestsChanged = true
		}
		statsMu.Unlock()
		if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
			return ErrNoTestsFound
		}
		return err
	}

	var runErr error
	if len(remainArgs) == 1 {
		runErr = processSingleArg(remainArgs[0], opts, runFn)
	} else {
		runErr = processMultiArg(remainArgs, opts, runFn)
	}

	if len(stats.Skipped) > 0 {
		runnerbuild.PrintSkippedSummary(stats.Skipped)
	}
	if stats.Total > 0 {
		stats.Elapsed = time.Since(start)
		runnerbuild.PrintResultSummary(opts, stats)
	}

	// Close metrics with run_end after summary.
	if rec != nil {
		_ = mergeNestSinkIntoRecorder(rec, rec.nestSink)
		warns := []string{}
		if metrics.ShouldWarnDefaultSuiteSlow(defaultSuite, stats.Total, stats.Elapsed, metrics.DefaultSuiteWarnThreshold) {
			warns = append(warns, "default_suite_slow")
		}
		exitOK := runErr == nil && stats.Passed >= stats.Total && stats.Total > 0
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
		if stats.NoTestsChanged || len(stats.Skipped) > 0 {
			return nil
		}
		return ErrNoTestsFound
	}
	if stats.Passed < stats.Total {
		return fmt.Errorf("%d of %d tests passed", stats.Passed, stats.Total)
	}
	return nil
}

func VetArgs(args []string) error {
	return processArgs(args, "vet", parseVetOptions, func(dir string, opts core.Options) error {
		return validate.RunWithOptions(dir, opts)
	})
}

func processArgs(args []string, cmdName string, parseFn func([]string) (core.Options, []string, error), processDirFn func(string, core.Options) error) error {
	opts, remainArgs, err := parseFn(args)
	if err != nil {
		return err
	}
	if len(remainArgs) < 1 {
		return fmt.Errorf("%s requires <dir>", cmdName)
	}
	if len(remainArgs) == 1 {
		return processSingleArg(remainArgs[0], opts, processDirFn)
	}
	return processMultiArg(remainArgs, opts, processDirFn)
}

func processSingleArg(arg string, opts core.Options, fn func(string, core.Options) error) error {
	if arg == "..." {
		return fmt.Errorf("bare '...' pattern is not supported; use './...' or 'path/...' instead")
	}
	if path_resolve.IsDotDotDotPattern(arg) {
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
	remainArgs, err := lessflags.Bool("-v,--verbose", &opts.Verbose).
		Bool("--rm", &opts.RemoveTemp).
		String("--gen-dir", &opts.GenDir).
		Int("-count", &opts.Count).
		Duration("--timeout", &opts.Timeout).
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
		Bool("--experiment-ref-instead-of-inline", &opts.ExperimentRefInsteadOfInline).
		Bool("--experiment-unified-package-per-doctest-tree", &opts.ExperimentUnifiedPackagePerDoctestTree).
		Parse(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	if colorFlag {
		opts.Color = core.ColorAlways
	}
	if noColorFlag {
		opts.Color = core.ColorNever
	}
	if opts.LabelAll && len(labelExprs) > 0 {
		return core.Options{}, nil, fmt.Errorf("--label-all and --label are mutually exclusive")
	}
	opts.LabelExprs = labelExprs
	// Unified package-per-tree implies ref-instead-of-inline generation.
	if opts.ExperimentUnifiedPackagePerDoctestTree {
		opts.ExperimentRefInsteadOfInline = true
	}

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
	if err := os.MkdirAll(gen, 0o755); err != nil {
		return fmt.Errorf("cold-cache: create gen dir %s: %w", gen, err)
	}
	opts.GenDir = gen

	gocacheTemp, err := os.MkdirTemp("", "doctest-cold-gocache-*")
	if err != nil {
		return fmt.Errorf("cold-cache: create isolated GOCACHE: %w", err)
	}
	opts.GoCache = gocacheTemp
	// So nested child processes that inherit the environment also see cold GOCACHE.
	_ = os.Setenv("GOCACHE", gocacheTemp)

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
