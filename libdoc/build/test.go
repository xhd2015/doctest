package build

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/leafcache"
	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

func Test(dir string, opts core.Options) error {
	_, err := TestWithStats(dir, opts)
	return err
}

func TestWithStats(dir string, opts core.Options) (TestRunStats, error) {
	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}

	absRoot, _ := filepath.Abs(dir)
	var phases []PhaseTiming
	track := func(name string, since time.Time) {
		phases = append(phases, PhaseTiming{Name: name, ElapsedNs: time.Since(since).Nanoseconds()})
	}

	// --- discover ---
	// SelectThenDeep: light scan (labels only) → filter → deep-parse run set.
	// Skipped labeled leaves are not fully parsed (default discovery speed).
	tDiscover := time.Now()
	allCases, err := core.DiscoverTreeCasesLight(dir)
	if err != nil {
		return TestRunStats{Phases: phases}, err
	}
	if opts.SubDir != "" {
		allCases = core.FilterBySubDir(allCases, dir, opts.SubDir)
	}

	var changedInfo core.ChangedRunInfo
	cases := allCases
	if opts.ChangedOnly {
		gitRoot, changedFiles, err := core.ChangedGitFiles(dir)
		if err != nil {
			track("discover", tDiscover)
			return TestRunStats{Phases: phases}, err
		}
		changedInfo = core.ChangedRunInfoForTree(allCases, dir, gitRoot, changedFiles)
		cases = core.FilterByChangedFiles(allCases, dir, gitRoot, changedFiles)
		if core.ShouldAnnounceChangedRun(changedInfo, opts.Verbose) {
			fmt.Fprintln(w, core.FormatDoctestAnnouncement(pathfmt.Short(dir), changedInfo, true, 0))
		}
		if len(cases) == 0 {
			track("discover", tDiscover)
			return TestRunStats{NoTestsChanged: true, Phases: phases}, nil
		}
	}

	cases, skipped := core.FilterCasesByLabel(cases, opts)
	for i := range skipped {
		skipped[i].DisplayPath = SkippedDisplayPath(dir, skipped[i].Path)
	}
	if len(cases) > 0 {
		cases, err = core.HydrateTreeCases(dir, cases)
		if err != nil {
			track("discover", tDiscover)
			return TestRunStats{Phases: phases, Skipped: skipped}, err
		}
	}
	track("discover", tDiscover)
	if len(cases) == 0 {
		if len(skipped) > 0 {
			stats := TestRunStats{Skipped: skipped, Phases: phases, Cases: nil}
			if !opts.SuppressResultSummary {
				sw := opts.Stdout
				if sw == nil {
					sw = os.Stdout
				}
				PrintSkippedSummaryTo(sw, skipped, opts.Verbose)
			}
			return stats, nil
		}
		return TestRunStats{Phases: phases}, fmt.Errorf("%s: no runnable test cases found", dir)
	}

	stats := TestRunStats{Total: len(cases), Skipped: skipped, Cases: cases}

	// --- materialize ---
	tMat := time.Now()
	ctx, err := newGenerateContext(dir, opts, cases, w, false, opts.Verbose)
	if err != nil {
		stats.Phases = phases
		return stats, err
	}
	track("materialize", tMat)
	ctx.installInterruptCleanup()
	defer ctx.Close()

	if opts.Verbose {
		ctx.announceRoots()
		if opts.ChangedOnly {
			fmt.Fprintln(w)
		} else {
			fmt.Fprintf(w, "doctest: %s\n\n", pathfmt.Short(dir))
		}
		if _, err := core.DiscoverTreeCasesVerbose(dir, w); err != nil {
			stats.Phases = phases
			return stats, err
		}
		fmt.Fprintf(w, "─── %d test cases\n\n", len(cases))
	} else if !opts.ChangedOnly {
		fmt.Fprintf(w, "doctest: %s (%d tests)\n", pathfmt.Short(dir), len(cases))
	}

	// --- generate ---
	tGen := time.Now()
	if err := ctx.writeCases(cases, false); err != nil {
		track("generate", tGen)
		stats.Phases = phases
		return stats, err
	}
	track("generate", tGen)

	// Multi-root workspace prep: generate only; caller runs __workspace/suite.
	if opts.GenerateOnly {
		stats.Phases = phases
		stats.GenRoot = ctx.genRoot
		stats.TreeRel = ctx.treeRel()
		stats.Unified = ctx.unifiedMode
		stats.AbsRoot = absRoot
		return stats, nil
	}

	// DOCTEST_DEBUG bypass-go-test: stop after generate (single-tree path).
	if opts.BypassGoTest {
		stats.Phases = phases
		stats.GenRoot = ctx.genRoot
		stats.TreeRel = ctx.treeRel()
		stats.Unified = ctx.unifiedMode
		stats.AbsRoot = absRoot
		stats.GoTestBypassed = true
		// Planned leaves stay in Total; do not claim pass.
		stats.Passed = 0
		return stats, nil
	}

	runDir, isSingleLeaf := ctx.runDir(absRoot, opts, cases)

	var packageArgs []string
	if ctx.unifiedMode {
		// Always suite-only packaging (even for a single leaf).
		runDir = ctx.scopedMultiRunDir(absRoot)
		var pkgErr error
		packageArgs, pkgErr = ctx.packageArgsForCases(runDir, absRoot, cases)
		if pkgErr != nil {
			stats.Phases = phases
			return stats, pkgErr
		}
	} else if isSingleLeaf {
		packageArgs = []string{"."}
	} else {
		runDir = ctx.scopedMultiRunDir(absRoot)
		var pkgErr error
		packageArgs, pkgErr = ctx.packageArgsForCases(runDir, absRoot, cases)
		if pkgErr != nil {
			stats.Phases = phases
			return stats, pkgErr
		}
	}

	// Flags only — packages are appended per go-test invocation (and per shard).
	flagArgs := []string{"test", "-mod=mod"}
	flagArgs = append(flagArgs, ctx.goCommandExtraArgs()...)
	if opts.Verbose {
		flagArgs = append(flagArgs, "-v")
	}
	if NeedsBuildVCSFlag(runDir) {
		flagArgs = append(flagArgs, "-buildvcs=false")
	}
	if opts.Count > 0 {
		flagArgs = append(flagArgs, fmt.Sprintf("-count=%d", opts.Count))
	}
	if opts.ForceWithFlagA {
		// go build/test -a: force rebuilding packages that are already up-to-date.
		flagArgs = append(flagArgs, "-a")
	}
	// nil = omit (go default 10m); non-nil including 0 = pass -timeout=…
	if opts.Timeout != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-timeout=%s", *opts.Timeout))
	}
	if opts.CPUProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-cpuprofile=%s", opts.CPUProfile))
	}
	if opts.MemProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-memprofile=%s", opts.MemProfile))
	}
	if opts.MemProfileRate != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-memprofilerate=%d", *opts.MemProfileRate))
	}
	if opts.BlockProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-blockprofile=%s", opts.BlockProfile))
	}
	if opts.BlockProfileRate != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-blockprofilerate=%d", *opts.BlockProfileRate))
	}
	if opts.MutexProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-mutexprofile=%s", opts.MutexProfile))
	}
	if opts.MutexProfileFraction != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-mutexprofilefraction=%d", *opts.MutexProfileFraction))
	}
	if opts.Trace != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-trace=%s", opts.Trace))
	}
	if opts.OutputDir != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-outputdir=%s", opts.OutputDir))
	}
	if opts.CoverProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-coverprofile=%s", opts.CoverProfile))
	}
	if opts.Cover {
		flagArgs = append(flagArgs, "-cover")
	}
	if opts.Race {
		flagArgs = append(flagArgs, "-race")
	}

	displayArgs := displayGoArgs(append(append([]string(nil), flagArgs...), packageArgs...))
	if opts.Verbose {
		fmt.Fprintf(w, "cd %s && go %s\n\n", pathfmt.Short(runDir), strings.Join(displayArgs, " "))
	} else {
		fmt.Fprintf(w, "cd %s && go %s\n", pathfmt.Short(runDir), strings.Join(displayArgs, " "))
	}

	sessionID := core.SessionIDFromOpts(opts)
	goCache := opts.GoCache
	// Prefer explicit nest sink on opts; else inherit process (suite child).
	if opts.MetricsNestSink == "" {
		if v := os.Getenv(envMetricsNestSink); v != "" {
			opts.MetricsNestSink = v
		}
	}

	// Programmatic leaf-cache: warm GetPass hits skip leaf bodies via suite env.
	// leafKeys + skipPaths also feed the JSON consumer for stream PutPass and
	// grey warm-skip progress dots (quiet + color).
	leafKeys, skipPaths := prepareLeafCache(absRoot, cases, opts)
	leafSkipEnv := leafcache.FormatSkipPaths(skipPaths)

	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	// go test rejects most profile/coverprofile/trace flags when multiple packages
	// are listed on one invocation. Run packages one-at-a-time when those flags are set.
	singlePkgInvocations := profileFlagsNeedSinglePackage(opts) && len(packageArgs) > 1

	// --- go_test ---
	// Always use go test -json for Pass/Fail/Run suite accounting. Verbose is
	// presentation only (stream more Output events); same counts as quiet.
	tGo := time.Now()
	style := newColorStyle(opts.Color, stdout)
	var result goTestJSONResult
	var goTestErr error
	if singlePkgInvocations {
		result, goTestErr = runGoTestJSONPerPackage(runDir, flagArgs, packageArgs, sessionID, goCache, opts.MetricsNestSink, leafSkipEnv, leafKeys, stdout, style, opts.Verbose)
	} else {
		result, goTestErr = runGoTestJSONShards(runDir, flagArgs, packageArgs, sessionID, goCache, opts.MetricsNestSink, leafSkipEnv, leafKeys, stdout, style, opts.Verbose)
	}
	goTestElapsed := time.Since(tGo)
	track("go_test", tGo)
	// Discovery planned count before Total is rewritten to actual_run.
	if stats.Planned == 0 {
		stats.Planned = stats.Total
	}
	if result.timeoutError != "" {
		stats.TimedOut = true
	}
	// Prefer JSON suite-leaf accounting when available. actual_run = pass+fail
	// (exclude runtime t.Skip from denominator). SkipCount is separate.
	// On timeout: never invent phantom passes from planned − failCount.
	actualRun := result.passCount + result.failCount
	if stats.TimedOut {
		stats.Passed = result.passCount
		stats.Total = actualRun
		stats.SkipCount = result.skipCount
	} else if actualRun > 0 || result.skipCount > 0 {
		stats.Passed = result.passCount
		stats.Total = actualRun
		stats.SkipCount = result.skipCount
	} else {
		stats.Passed = passedCases(stats.Total, result.failCount)
	}
	if ctx.unifiedMode {
		stats.LeafTimings = leafTimingsFromSubtests(cases, result, goTestElapsed)
	} else {
		stats.LeafTimings = leafTimingsFromPackages(cases, packageArgs, isSingleLeaf, result, goTestElapsed)
	}

	// Summary Cached is leaf-cache-only:
	//   - Cached = number of warm leaf-cache skips (GetPass hits used for skip)
	//   - full go package (cached) expands to N only when every leaf key also
	//     hits (otherwise go DCE/testcache can stay warm while spine text
	//     changed and leaf keys miss — product wants 0 Cached then)
	//   - leaf-cache bypass (-count / -a / --no-leaf-cache): always 0
	result.cachedCount = leafCachedSummary(len(cases), skipPaths, result.anyPackageCached(),
		leafcache.SkipEnabled(opts.Count, opts.ForceWithFlagA, opts.NoLeafCache))
	recordLeafCachePasses(leafKeys, result.suiteLeafFailed, goTestErr == nil && result.failCount == 0)

	// Quiet path: compact progress summary. Verbose already streamed Output events.
	// Progress stays finished-only (no Cancelled segment).
	if !opts.Verbose {
		fmt.Fprintln(stdout, formatSummary(style, result.passCount+result.failCount, result.passCount, result.failCount, result.cachedCount, goTestElapsed))
	}

	// Print order: progress → fail dumps → Error/hint → PASS/FAIL.
	// Quiet path buffers fail details for post-run print. Verbose streamed live.
	if !opts.Verbose {
		for _, line := range result.failLines {
			fmt.Fprintln(stdout, line)
		}
		for _, line := range result.detailLines {
			fmt.Fprintln(stdout, line)
		}
	}
	if len(result.stderrData) > 0 {
		stdout.Write(result.stderrData)
	}
	// Timeout Error/hint on stdout (before FAIL) so user-facing order is correct.
	printGoTestTimeoutError(stdout, result, style)

	if !opts.SuppressResultSummary {
		stats.Elapsed = goTestElapsed
		PrintSkippedSummaryTo(stdout, stats.Skipped, opts.Verbose)
		PrintResultSummary(opts, stats)
	}

	if goTestErr != nil {
		stats.Phases = phases
		return stats, fmt.Errorf("go test: %w", goTestErr)
	}

	tPost := time.Now()
	if err := ctx.syncDump(); err != nil {
		track("post", tPost)
		stats.Phases = phases
		return stats, err
	}
	if d := time.Since(tPost); d > time.Millisecond {
		track("post", tPost)
	}
	if !opts.Verbose {
		fmt.Fprintln(w)
	}
	stats.Phases = phases
	return stats, nil
}

// UnifiedSuiteTestName is the go test function that iterates registered leaves.
const UnifiedSuiteTestName = "TestDoctestSuite"

// leafTimingsFromSubtests attributes go test -json subtest Elapsed under the
// unified suite (t.Run(encodedPath, …) → TestDoctestSuite/<path with / → __>).
// Workspace nested names are TestDoctestSuite/<tree>/<leaf>; the leaf segment
// (after the last "/") is used to match tree-relative case paths.
// When unmappable, elapsed stays 0 (do not clone suite wall onto every leaf).
func leafTimingsFromSubtests(cases []core.TreeCase, result goTestJSONResult, goTestWall time.Duration) []LeafTiming {
	out := make([]LeafTiming, 0, len(cases))
	byRel := map[string]int64{}
	prefix := UnifiedSuiteTestName + "/"
	for test, ns := range result.testElapsedNs {
		if !strings.HasPrefix(test, prefix) {
			continue
		}
		rel := test[len(prefix):]
		if rel == "" {
			continue
		}
		// Workspace: tree/leaf → keep leaf segment for case.Path match.
		if i := strings.LastIndex(rel, "/"); i >= 0 {
			rel = rel[i+1:]
		}
		// Decode suite subtest encoding: "/" was replaced with "__".
		rel = strings.ReplaceAll(rel, "__", "/")
		// Prefer the deepest subtest elapsed (leaf path); keep max if duplicates.
		if ns > byRel[rel] {
			byRel[rel] = ns
		}
	}
	// Single-leaf unified: if only package elapsed exists, fall back once.
	if len(cases) == 1 && len(byRel) == 0 {
		var ns int64
		for _, e := range result.pkgElapsedNs {
			if e > ns {
				ns = e
			}
		}
		if ns == 0 {
			ns = goTestWall.Nanoseconds()
		}
		return []LeafTiming{{Path: cases[0].Path, ElapsedNs: ns}}
	}
	for _, tc := range cases {
		lt := LeafTiming{Path: tc.Path}
		slash := filepath.ToSlash(tc.Path)
		if ns, ok := byRel[slash]; ok {
			lt.ElapsedNs = ns
		} else if ns, ok := byRel[tc.Path]; ok {
			lt.ElapsedNs = ns
		} else {
			// Match by suffix (absolute vs relative discovery paths).
			for rel, ns := range byRel {
				if slash == rel || strings.HasSuffix(slash, "/"+rel) || strings.HasSuffix(rel, "/"+slash) {
					lt.ElapsedNs = ns
					break
				}
			}
		}
		out = append(out, lt)
	}
	return out
}

// leafTimingsFromPackages attributes go test -json package Elapsed to leaves.
// When unmappable, multi-leaf elapsed stays 0 (do not clone tree wall onto every leaf).
func leafTimingsFromPackages(cases []core.TreeCase, packageArgs []string, isSingleLeaf bool, result goTestJSONResult, goTestWall time.Duration) []LeafTiming {
	out := make([]LeafTiming, 0, len(cases))
	if isSingleLeaf && len(cases) == 1 {
		cached := false
		for _, c := range result.pkgCached {
			if c {
				cached = true
				break
			}
		}
		// Prefer package elapsed when present.
		var ns int64
		for _, e := range result.pkgElapsedNs {
			if e > ns {
				ns = e
			}
		}
		if ns == 0 {
			ns = goTestWall.Nanoseconds()
		}
		return []LeafTiming{{Path: cases[0].Path, ElapsedNs: ns, Cached: cached}}
	}
	for _, tc := range cases {
		lt := LeafTiming{Path: tc.Path}
		slash := filepath.ToSlash(tc.Path)
		for pkg, ns := range result.pkgElapsedNs {
			if packageMatchesLeaf(pkg, slash) {
				lt.ElapsedNs = ns
				lt.Cached = result.pkgCached[pkg]
				break
			}
		}
		// Fallback: packageArgs "./rel" vs leaf path suffix.
		if lt.ElapsedNs == 0 {
			for _, arg := range packageArgs {
				rel := strings.TrimPrefix(filepath.ToSlash(arg), "./")
				if rel == slash || strings.HasSuffix(rel, "/"+slash) || strings.HasSuffix(slash, "/"+rel) || rel == filepath.Base(slash) {
					// try match any pkg ending with rel
					for pkg, ns := range result.pkgElapsedNs {
						if strings.HasSuffix(filepath.ToSlash(pkg), "/"+rel) || strings.HasSuffix(filepath.ToSlash(pkg), rel) {
							lt.ElapsedNs = ns
							lt.Cached = result.pkgCached[pkg]
							break
						}
					}
				}
				if lt.ElapsedNs > 0 {
					break
				}
			}
		}
		out = append(out, lt)
	}
	return out
}

func packageMatchesLeaf(pkg, leafSlash string) bool {
	p := filepath.ToSlash(pkg)
	if p == leafSlash {
		return true
	}
	if strings.HasSuffix(p, "/"+leafSlash) {
		return true
	}
	// last path segment(s)
	base := filepath.Base(leafSlash)
	return strings.HasSuffix(p, "/"+base) && strings.Contains(p, leafSlash)
}

type goTestJSONResult struct {
	passCount   int
	failCount   int
	skipCount   int // suite-leaf Action "skip" (t.Skip); not package-level
	cachedCount int
	failLines   []string
	detailLines []string
	stderrData  []byte
	// timeoutError is a clear user-facing line when go test panics with
	// "test timed out after <d>" (JSON Output events are otherwise buffered
	// under the test name and dropped because timeout emits no per-test fail).
	timeoutError  string
	pkgElapsedNs  map[string]int64 // import path -> package-level Elapsed from -json
	pkgCached     map[string]bool
	testElapsedNs map[string]int64 // full test name (incl. subtests) -> Elapsed
	// suiteLeafFailed maps tree-relative leaf paths (slash form) that failed.
	suiteLeafFailed map[string]bool
}

func (r goTestJSONResult) anyPackageCached() bool {
	for _, c := range r.pkgCached {
		if c {
			return true
		}
	}
	return false
}

func mergeGoTestJSONResult(dst *goTestJSONResult, src goTestJSONResult) {
	dst.passCount += src.passCount
	dst.failCount += src.failCount
	dst.skipCount += src.skipCount
	dst.cachedCount += src.cachedCount
	dst.failLines = append(dst.failLines, src.failLines...)
	dst.detailLines = append(dst.detailLines, src.detailLines...)
	dst.stderrData = append(dst.stderrData, src.stderrData...)
	if dst.timeoutError == "" && src.timeoutError != "" {
		dst.timeoutError = src.timeoutError
	}
	if src.pkgElapsedNs != nil {
		if dst.pkgElapsedNs == nil {
			dst.pkgElapsedNs = make(map[string]int64)
		}
		for k, v := range src.pkgElapsedNs {
			dst.pkgElapsedNs[k] = v
		}
	}
	if src.pkgCached != nil {
		if dst.pkgCached == nil {
			dst.pkgCached = make(map[string]bool)
		}
		for k, v := range src.pkgCached {
			dst.pkgCached[k] = v
		}
	}
	if src.suiteLeafFailed != nil {
		if dst.suiteLeafFailed == nil {
			dst.suiteLeafFailed = make(map[string]bool)
		}
		for k, v := range src.suiteLeafFailed {
			if v {
				dst.suiteLeafFailed[k] = true
			}
		}
	}
	if src.testElapsedNs != nil {
		if dst.testElapsedNs == nil {
			dst.testElapsedNs = make(map[string]int64)
		}
		for k, v := range src.testElapsedNs {
			dst.testElapsedNs[k] = v
		}
	}
}

// packageTestShards splits package paths across workers for concurrent go test
// processes (go test itself runs listed packages serially).
func packageTestShards(pkgs []string, workers int) [][]string {
	if len(pkgs) == 0 {
		return nil
	}
	if workers <= 1 || len(pkgs) == 1 {
		return [][]string{pkgs}
	}
	if workers > len(pkgs) {
		workers = len(pkgs)
	}
	shards := make([][]string, workers)
	for i, p := range pkgs {
		shards[i%workers] = append(shards[i%workers], p)
	}
	out := shards[:0]
	for _, s := range shards {
		if len(s) > 0 {
			out = append(out, s)
		}
	}
	return out
}

// profileFlagsNeedSinglePackage reports whether go test would reject the
// configured profile/cover/trace flags when multiple packages are listed.
func profileFlagsNeedSinglePackage(opts core.Options) bool {
	return opts.CPUProfile != "" ||
		opts.MemProfile != "" ||
		opts.BlockProfile != "" ||
		opts.MutexProfile != "" ||
		opts.Trace != "" ||
		opts.CoverProfile != "" ||
		opts.OutputDir != ""
}

// envMetricsNestSink matches metrics.EnvMetricsNestSink (avoid build→metrics import).
const envMetricsNestSink = "DOCTEST_METRICS_NEST_SINK"

// goTestEnv builds the environment for a child `go test` process, including
// DOCTEST_SESSION_ID, optional isolated GOCACHE, and nest sink when present.
// Nest sink is process-lifetime for the suite binary (leaves copy into
// d.Metrics.NestSink via read-only getenv); prefer opts.MetricsNestSink +
// goTestEnvFromOpts over mid-suite process env mutation.
func goTestEnv(sessionID, goCache string) []string {
	return goTestEnvFull(sessionID, goCache, "", "", "")
}

// goTestEnvFull builds child env with key-replace (core.ChildEnv) so SESSION_ID,
// GOCACHE, nest sink, GOWORK, and leaf-skip paths override process values
// without blind append duplicates.
func goTestEnvFull(sessionID, goCache, nestSink, goWork, leafSkipPaths string) []string {
	overrides := []string{core.DoctestSessionIDEnv + "=" + sessionID}
	if goCache != "" {
		overrides = append(overrides, "GOCACHE="+goCache)
	}
	if nestSink != "" {
		overrides = append(overrides, envMetricsNestSink+"="+nestSink)
	} else if v := os.Getenv(envMetricsNestSink); v != "" {
		// Inherit outer nest sink for nested go test (read-only).
		overrides = append(overrides, envMetricsNestSink+"="+v)
	}
	if goWork != "" {
		// Multi-module workspace hub (go.work). Empty means use process default / off.
		overrides = append(overrides, "GOWORK="+goWork)
	}
	if leafSkipPaths != "" {
		overrides = append(overrides, leafcache.EnvSkipPaths+"="+leafSkipPaths)
	}
	return core.ChildEnv(nil, overrides...)
}

func goTestEnvFromOpts(sessionID string, opts core.Options) []string {
	if sessionID == "" {
		sessionID = core.SessionIDFromOpts(opts)
	}
	return goTestEnvFull(sessionID, opts.GoCache, opts.MetricsNestSink, "", "")
}

// runGoTestJSONPerPackage runs one go test process per package (serial) so
// profile flags that go rejects with multi-package lists still work.
func runGoTestJSONPerPackage(runDir string, flagArgs, packageArgs []string, sessionID, goCache, nestSink, leafSkipPaths string, leafKeys map[string]string, stdout io.Writer, style colorStyle, verbose bool) (goTestJSONResult, error) {
	var merged goTestJSONResult
	var firstErr error
	for _, pkg := range packageArgs {
		args := append(append([]string(nil), flagArgs...), pkg)
		res, err := runGoTestJSONOnce(runDir, args, sessionID, goCache, nestSink, "", leafSkipPaths, leafKeys, stdout, style, verbose)
		mergeGoTestJSONResult(&merged, res)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return merged, firstErr
}

// goTestSlots bounds concurrent `go test` processes in this doctest process.
// Nested self-test subprocesses have their own pool (separate process).
var goTestSlots = make(chan struct{}, 4)

// goTestModRootMu serializes go test under the same module root so parallel
// ./... trees sharing one gen root (e.g. cold mapping-gen-cold) do not race.
var goTestModRootMu sync.Map // absModDir -> *sync.Mutex

func lockGoTestModule(runDir string) func() {
	modDir := runDir
	for {
		if _, err := os.Stat(filepath.Join(modDir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(modDir)
		if parent == modDir {
			modDir = runDir
			break
		}
		modDir = parent
	}
	abs, err := filepath.Abs(modDir)
	if err != nil {
		abs = modDir
	}
	v, _ := goTestModRootMu.LoadOrStore(abs, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

func runGoTestJSONShards(runDir string, flagArgs, packageArgs []string, sessionID, goCache, nestSink, leafSkipPaths string, leafKeys map[string]string, stdout io.Writer, style colorStyle, verbose bool) (goTestJSONResult, error) {
	// Single go test process per tree. Package sharding multiplies nested
	// self-test fan-out and has raced go.mod; wall cut is tree concurrency +
	// suite-level t.Parallel.
	workers := 1
	shards := packageTestShards(packageArgs, workers)
	if len(shards) <= 1 {
		args := append(append([]string(nil), flagArgs...), packageArgs...)
		return runGoTestJSONOnce(runDir, args, sessionID, goCache, nestSink, "", leafSkipPaths, leafKeys, stdout, style, verbose)
	}

	// Multi-shard: readonly module mode so concurrent go tests share genDir safely.
	shardFlags := make([]string, 0, len(flagArgs))
	for _, a := range flagArgs {
		if a == "-mod=mod" {
			shardFlags = append(shardFlags, "-mod=readonly")
			continue
		}
		shardFlags = append(shardFlags, a)
	}

	var (
		mu       sync.Mutex
		merged   goTestJSONResult
		firstErr error
		wg       sync.WaitGroup
	)
	wg.Add(len(shards))
	for _, shard := range shards {
		shard := shard
		go func() {
			defer wg.Done()
			args := append(append([]string(nil), shardFlags...), shard...)
			// Locked stdout keeps progress dots incremental and non-interleaved by byte.
			res, err := runGoTestJSONOnce(runDir, args, sessionID, goCache, nestSink, "", leafSkipPaths, leafKeys, &lockedWriter{w: stdout, mu: &mu}, style, verbose)
			mu.Lock()
			defer mu.Unlock()
			mergeGoTestJSONResult(&merged, res)
			if err != nil && firstErr == nil {
				firstErr = err
			}
		}()
	}
	wg.Wait()
	return merged, firstErr
}

type lockedWriter struct {
	w  io.Writer
	mu *sync.Mutex
}

func (l *lockedWriter) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.w.Write(p)
}

func runGoTestJSONOnce(runDir string, testArgs []string, sessionID, goCache, nestSink, goWork, leafSkipPaths string, leafKeys map[string]string, stdout io.Writer, style colorStyle, verbose bool) (goTestJSONResult, error) {
	goTestSlots <- struct{}{}
	defer func() { <-goTestSlots }()

	unlockMod := lockGoTestModule(runDir)
	defer unlockMod()

	// Always -json for suite accounting. Drop -v from the real invocation:
	// go test -json already emits Output for fmt.Print/t.Logf; combining -v
	// makes test2json re-parse framing lines and invent phantom fail events
	// when nested suites print "--- FAIL:" / "=== RUN" into a passing leaf.
	// Display still shows -v (flagArgs) so verbose-go-flag / user-facing cd
	// lines keep advertising presentation mode.
	execArgs := make([]string, 0, len(testArgs)+1)
	for _, a := range testArgs {
		if a == "-v" {
			continue
		}
		execArgs = append(execArgs, a)
	}
	execArgs = append(execArgs, "-json")
	goTestCmd := exec.Command("go", execArgs...)
	goTestCmd.Dir = runDir
	goTestCmd.Env = goTestEnvFull(sessionID, goCache, nestSink, goWork, leafSkipPaths)

	stdoutPipe, err := goTestCmd.StdoutPipe()
	if err != nil {
		return goTestJSONResult{}, fmt.Errorf("stdout pipe: %w", err)
	}
	stderrPipe, err := goTestCmd.StderrPipe()
	if err != nil {
		return goTestJSONResult{}, fmt.Errorf("stderr pipe: %w", err)
	}
	if err := goTestCmd.Start(); err != nil {
		return goTestJSONResult{}, fmt.Errorf("go test start: %w", err)
	}

	// Stream PutPass: open store once for this consumer so mid-run Ctrl-C still
	// leaves already-passed leaves durable. Fail never PutPass (only pass path).
	var streamStore *leafcache.Store
	if len(leafKeys) > 0 {
		if root, err := leafcache.ResolveStoreRoot(); err == nil {
			streamStore, _ = leafcache.NewStore(root)
		}
	}
	// This-run warm skip set (bare paths and/or FormatLeafIdentityEnv tokens).
	warmSkip := leafcache.ParseSkipPaths(leafSkipPaths)

	var res goTestJSONResult
	res.pkgElapsedNs = make(map[string]int64)
	res.pkgCached = make(map[string]bool)
	res.testElapsedNs = make(map[string]int64)
	res.suiteLeafFailed = make(map[string]bool)
	var stdoutWg sync.WaitGroup
	stdoutWg.Add(1)
	go func() {
		defer stdoutWg.Done()
		type goTestEvent struct {
			Action  string  `json:"Action"`
			Package string  `json:"Package"`
			Test    string  `json:"Test"`
			Output  string  `json:"Output"`
			Elapsed float64 `json:"Elapsed"`
		}
		packageCached := make(map[string]bool)
		failedTests := make(map[string]bool)
		testOutputs := make(map[string][]string)
		// suiteLeaf* count progress + summary from unified suite subtests
		// (TestDoctestSuite/<leafPath>). When any leaf subtest is seen, skip
		// package-level dots and use leaf pass/fail/skip/cached for the summary
		// so one suite binary still reports one result per leaf (not 1 package).
		suiteLeafPass := 0
		suiteLeafFail := 0
		suiteLeafSkip := 0
		// Dedup suite leaves by package+name: multi-module go.work runs multiple
		// suite packages that often share leaf subtest names (e.g. "simple").
		// Keying only on test name collapsed those to one count.
		suiteLeafSeen := make(map[string]bool)
		testKey := func(pkg, test string) string { return pkg + "\x00" + test }
		// Suite leaf progress: per-tree suite uses flat TestDoctestSuite/<leaf>.
		// Workspace suite nests TestDoctestSuite/<tree>/<leaf> — count only leaves
		// (two+ segments), not the tree parallel parent.
		isCountableSuiteLeaf := func(pkg, test string) bool {
			prefix := UnifiedSuiteTestName + "/"
			if !strings.HasPrefix(test, prefix) {
				return false
			}
			rest := test[len(prefix):]
			if rest == "" {
				return false
			}
			// Nested suites (workspace multi-tree or multi-mod hub): only count
			// leaves (TestDoctestSuite/<tree>/<leaf>), not the tree parent node.
			if strings.Contains(pkg, "/"+core.WorkspaceDirName+"/") ||
				strings.HasSuffix(pkg, "/"+core.WorkspaceDirName+"/"+core.WorkspaceSuiteDirName) ||
				strings.Contains(pkg, core.WorkspaceDirName+"/"+core.WorkspaceSuitePkgName) ||
				strings.Contains(pkg, "__workspace") ||
				strings.Contains(pkg, "/"+HubDirName+"/") ||
				strings.HasPrefix(pkg, hubModulePath+"/") ||
				pkg == hubModulePath+"/suite" {
				return strings.Contains(rest, "/")
			}
			return true
		}
		// outcome: "pass", "fail", or "skip"
		countSuiteLeaf := func(pkg, test string, outcome string) {
			key := testKey(pkg, test)
			if suiteLeafSeen[key] {
				return
			}
			suiteLeafSeen[key] = true
			switch outcome {
			case "fail":
				suiteLeafFail++
				if rel := suiteLeafRelPath(test); rel != "" {
					res.suiteLeafFailed[rel] = true
				}
			case "skip":
				suiteLeafSkip++
			default:
				suiteLeafPass++
			}
		}
		flushTestOutput := func(key string) {
			if failedTests[key] {
				return
			}
			failedTests[key] = true
			if buf := testOutputs[key]; len(buf) > 0 {
				res.detailLines = append(res.detailLines, buf...)
				delete(testOutputs, key)
			}
		}
		recordPkg := func(pkg string, elapsed float64, cached bool) {
			if pkg == "" {
				return
			}
			if elapsed > 0 {
				res.pkgElapsedNs[pkg] = int64(elapsed * float64(time.Second))
			}
			if cached {
				res.pkgCached[pkg] = true
				packageCached[pkg] = true
			}
		}
		recordTest := func(test string, elapsed float64) {
			if test == "" {
				return
			}
			// Keep max elapsed if the same name appears more than once.
			ns := int64(elapsed * float64(time.Second))
			if ns > res.testElapsedNs[test] {
				res.testElapsedNs[test] = ns
			}
		}
		// Stream PutPass for a countable suite-leaf pass (best-effort).
		// Fail never calls this. End-of-run recordLeafCachePasses remains as
		// idempotent reconcile for clean completions.
		streamPutPassLeaf := func(leafRel string) {
			if streamStore == nil || leafRel == "" || len(leafKeys) == 0 {
				return
			}
			if key := streamKeyForLeaf(leafKeys, leafRel); key != "" {
				_ = streamStore.PutPass(key)
			}
		}
		// Quiet progress dots only; verbose streams Output events instead.
		// Warm leaf-cache skips (identity in this-run skip set) print grey
		// when color is on; executed passes stay plain "."; fails stay red.
		writePassDot := func(leafRel string) {
			if verbose {
				return
			}
			if leafRel != "" && style.enabled && isWarmSkipLeaf(warmSkip, leafRel) {
				fmt.Fprint(stdout, style.gray("."))
				return
			}
			stdout.Write([]byte("."))
		}
		writeFailDot := func() {
			if !verbose {
				fmt.Fprint(stdout, style.red("."))
			}
		}
		decoder := json.NewDecoder(stdoutPipe)
		for {
			var ev goTestEvent
			if err := decoder.Decode(&ev); err != nil {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "read go test json: %v\n", err)
				}
				break
			}
			switch ev.Action {
			case "output":
				// Timeout panics are Output events under a Test name; go test then
				// only emits package-level fail (no per-test fail), so buffered
				// lines would be dropped — capture a clear Error line here.
				if res.timeoutError == "" {
					if msg := goTestTimeoutErrorLine(ev.Output); msg != "" {
						res.timeoutError = msg
					}
				}
				// Verbose: stream raw Output (presentation only; counts from pass/fail).
				if verbose && ev.Output != "" {
					stdout.Write([]byte(ev.Output))
				}
				trimmed := strings.TrimSpace(ev.Output)
				if ev.Test == "" && (strings.HasPrefix(trimmed, "ok ") || strings.HasPrefix(trimmed, "ok\t")) {
					if strings.Contains(ev.Output, "(cached)") {
						packageCached[ev.Package] = true
					}
				}
				// Quiet path buffers fail details for post-run print. Verbose already
				// streamed above — still track package FAIL lines for accounting paths.
				if !verbose {
					if ev.Test == "" && (strings.HasPrefix(trimmed, "FAIL\t") || strings.HasPrefix(trimmed, "FAIL ")) {
						res.failLines = append(res.failLines, trimmed)
					}
					if trimmed != "" && trimmed != "PASS" {
						if ev.Test != "" {
							key := testKey(ev.Package, ev.Test)
							line := strings.TrimRight(ev.Output, "\n")
							if strings.Contains(ev.Output, "--- FAIL:") {
								flushTestOutput(key)
								res.detailLines = append(res.detailLines, line)
							} else if failedTests[key] {
								res.detailLines = append(res.detailLines, line)
							} else {
								testOutputs[key] = append(testOutputs[key], line)
							}
						} else if !strings.HasPrefix(trimmed, "ok ") && !strings.HasPrefix(trimmed, "ok\t") &&
							!strings.HasPrefix(trimmed, "FAIL\t") && !strings.HasPrefix(trimmed, "FAIL ") {
							res.detailLines = append(res.detailLines, trimmed)
						}
					}
				}
			case "pass":
				if ev.Test != "" {
					recordTest(ev.Test, ev.Elapsed)
					delete(testOutputs, testKey(ev.Package, ev.Test))
					if isCountableSuiteLeaf(ev.Package, ev.Test) {
						before := suiteLeafPass + suiteLeafFail + suiteLeafSkip
						countSuiteLeaf(ev.Package, ev.Test, "pass")
						if suiteLeafPass+suiteLeafFail+suiteLeafSkip > before {
							leafRel := suiteLeafRelPath(ev.Test)
							// PutPass before the progress dot so interrupt-after-dots
							// harnesses (and Ctrl-C) observe a durable store entry
							// once a pass dot is visible.
							streamPutPassLeaf(leafRel)
							writePassDot(leafRel)
						}
					}
					continue
				}
				res.passCount++
				if packageCached[ev.Package] {
					res.cachedCount++
				}
				recordPkg(ev.Package, ev.Elapsed, packageCached[ev.Package])
				if suiteLeafPass+suiteLeafFail+suiteLeafSkip == 0 {
					writePassDot("")
				}
			case "fail":
				if ev.Test != "" {
					recordTest(ev.Test, ev.Elapsed)
					key := testKey(ev.Package, ev.Test)
					flushTestOutput(key)
					if isCountableSuiteLeaf(ev.Package, ev.Test) {
						before := suiteLeafPass + suiteLeafFail + suiteLeafSkip
						countSuiteLeaf(ev.Package, ev.Test, "fail")
						if suiteLeafPass+suiteLeafFail+suiteLeafSkip > before {
							writeFailDot()
						}
					}
					continue
				}
				res.failCount++
				recordPkg(ev.Package, ev.Elapsed, false)
				if suiteLeafPass+suiteLeafFail+suiteLeafSkip == 0 {
					writeFailDot()
				}
			case "skip":
				// Runtime t.Skip on suite leaves: count separately (not pass/fail).
				if ev.Test != "" {
					recordTest(ev.Test, ev.Elapsed)
					delete(testOutputs, testKey(ev.Package, ev.Test))
					if isCountableSuiteLeaf(ev.Package, ev.Test) {
						countSuiteLeaf(ev.Package, ev.Test, "skip")
					}
					continue
				}
			}
		}
		// Prefer leaf-level counts for unified suite so multi-leaf trees report
		// N Run / N Pass instead of a single package result.
		if suiteLeafPass+suiteLeafFail+suiteLeafSkip > 0 {
			res.passCount = suiteLeafPass
			res.failCount = suiteLeafFail
			res.skipCount = suiteLeafSkip
			anyCached := false
			for _, c := range packageCached {
				if c {
					anyCached = true
					break
				}
			}
			if anyCached {
				res.cachedCount = suiteLeafPass
			} else {
				res.cachedCount = 0
			}
		}
	}()

	stderrData, _ := io.ReadAll(stderrPipe)
	stdoutWg.Wait()
	err = goTestCmd.Wait()
	res.stderrData = stderrData
	// go test may also print the panic only on stderr in edge cases.
	if res.timeoutError == "" {
		if msg := goTestTimeoutErrorLine(string(stderrData)); msg != "" {
			res.timeoutError = msg
		}
	}
	// Avoid false suite-level timeout when nested fail dumps contain
	// "test timed out after" but this go test process succeeded (e.g. outer
	// harness leaf that asserts nested timeout messaging).
	if err == nil {
		res.timeoutError = ""
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

// goTestTimeoutErrorLine returns a clear user-facing timeout message when s
// contains go test's "test timed out after <duration>" panic phrase.
// Example:
//
//	Error: go test timed out after 2s
//	hint: increase with -timeout=DURATION (e.g. -timeout=30m; -timeout=0 disables)
func goTestTimeoutErrorLine(s string) string {
	const marker = "test timed out after "
	idx := strings.Index(s, marker)
	if idx < 0 {
		return ""
	}
	rest := strings.TrimSpace(s[idx+len(marker):])
	if rest == "" {
		return ""
	}
	dur := rest
	if i := strings.IndexAny(rest, " \t\r\n"); i >= 0 {
		dur = rest[:i]
	}
	if dur == "" {
		return ""
	}
	return "Error: go test timed out after " + dur +
		"\nhint: increase with -timeout=DURATION (e.g. -timeout=30m; -timeout=0 disables)"
}

// printGoTestTimeoutError writes the surfaced timeout Error + hint lines (if any)
// to stdout so they appear before the end FAIL summary on the user-facing stream.
// When color is enabled: Error line red, hint line gray.
func printGoTestTimeoutError(stdout io.Writer, result goTestJSONResult, style colorStyle) {
	if result.timeoutError == "" || stdout == nil {
		return
	}
	for _, line := range strings.Split(result.timeoutError, "\n") {
		if line == "" {
			continue
		}
		if style.enabled {
			switch {
			case strings.HasPrefix(line, "Error:"):
				line = style.red(line)
			case strings.HasPrefix(line, "hint:"):
				line = style.gray(line)
			}
		}
		fmt.Fprintln(stdout, line)
	}
}

// suiteLeafRelPath decodes a go test -json Test name under TestDoctestSuite/
// into a tree-relative leaf path (reversing RunAll's "/" → "__" encoding).
// Workspace nests as TestDoctestSuite/<tree>/<leaf>; the leaf segment is returned
// when multiple segments are present (tree-relative leaf path only).
func suiteLeafRelPath(test string) string {
	prefix := UnifiedSuiteTestName + "/"
	if !strings.HasPrefix(test, prefix) {
		return ""
	}
	rest := test[len(prefix):]
	if rest == "" {
		return ""
	}
	// Workspace: tree/leaf (real slash from nested t.Run) — keep leaf segment.
	if i := strings.LastIndex(rest, "/"); i >= 0 {
		rest = rest[i+1:]
	}
	return strings.ReplaceAll(rest, "__", "/")
}

// streamKeyForLeaf resolves a store key for a suite leaf relative path.
// Single-tree prepareLeafCache keys by bare path; workspace keys use
// FormatLeafIdentity (NUL-separated). When multiple FormatLeafIdentity keys
// share the same bare leaf (multi-tree same relpath), returns "" so stream
// PutPass skips — end-of-run RecordPasses still reconciles uniquely.
func streamKeyForLeaf(leafKeys map[string]string, leafRel string) string {
	if len(leafKeys) == 0 || leafRel == "" {
		return ""
	}
	if k, ok := leafKeys[leafRel]; ok {
		return k
	}
	slash := filepath.ToSlash(leafRel)
	if k, ok := leafKeys[slash]; ok {
		return k
	}
	// Unique FormatLeafIdentity suffix match (treeRoot\x00leafRel).
	var found string
	n := 0
	needle := "\x00" + leafRel
	needleSlash := "\x00" + slash
	for id, k := range leafKeys {
		if id == leafRel || id == slash ||
			strings.HasSuffix(id, needle) || strings.HasSuffix(id, needleSlash) {
			found = k
			n++
			if n > 1 {
				return ""
			}
		}
	}
	if n == 1 {
		return found
	}
	return ""
}

// isWarmSkipLeaf reports whether leafRel is in this-run warm skip set.
// Skip tokens are bare paths (single-tree) and/or FormatLeafIdentityEnv
// (absRoot\tleafRel for workspace).
func isWarmSkipLeaf(warmSkip map[string]struct{}, leafRel string) bool {
	if len(warmSkip) == 0 || leafRel == "" {
		return false
	}
	if _, ok := warmSkip[leafRel]; ok {
		return true
	}
	slash := filepath.ToSlash(leafRel)
	if _, ok := warmSkip[slash]; ok {
		return true
	}
	suffix := "\t" + leafRel
	suffixSlash := "\t" + slash
	for p := range warmSkip {
		if strings.HasSuffix(p, suffix) || strings.HasSuffix(p, suffixSlash) {
			return true
		}
	}
	return false
}

// leafCachedSummary computes summary Cached for the leaf-cache product.
// When skip is disabled (-count / -a / --no-leaf-cache), always 0.
// Otherwise Cached is the leaf-skip count; full go package (cached) expands to
// all N leaves only when every leaf is also a warm GetPass hit.
func leafCachedSummary(nCases int, skipPaths []string, anyPkgCached, skipEnabled bool) int {
	if !skipEnabled {
		return 0
	}
	nSkip := len(skipPaths)
	if anyPkgCached && nCases > 0 && nSkip == nCases {
		return nCases
	}
	return nSkip
}

// prepareLeafCache computes leaf keys and warm skip paths for this tree run.
// Uses leafcache.PreparePassPlan, then remaps tree-qualified identities back to
// bare tree-relative paths so the suite skip env (DOCTEST_LEAF_CACHE_SKIP_PATHS)
// and suiteLeafFailed maps stay path-based for single-tree runs.
// Store I/O errors are ignored (best-effort; suite continues).
func prepareLeafCache(treeRoot string, cases []core.TreeCase, opts core.Options) (keys map[string]string, skipPaths []string) {
	keys = make(map[string]string, len(cases))
	if len(cases) == 0 {
		return keys, nil
	}
	storeRoot, err := leafcache.ResolveStoreRoot()
	if err != nil {
		return keys, nil
	}
	store, err := leafcache.NewStore(storeRoot)
	if err != nil {
		return keys, nil
	}
	goVer := runtime.Version()
	enabled := leafcache.SkipEnabled(opts.Count, opts.ForceWithFlagA, opts.NoLeafCache)
	leaves := make([]leafcache.LeafRef, 0, len(cases))
	idToPath := make(map[string]string, len(cases))
	for _, tc := range cases {
		leaves = append(leaves, leafcache.LeafRef{TreeRoot: treeRoot, LeafRel: tc.Path})
		idToPath[leafcache.FormatLeafIdentity(treeRoot, tc.Path)] = tc.Path
	}
	plan, err := leafcache.PreparePassPlan(store, leaves, goVer, enabled)
	if err != nil {
		return keys, nil
	}
	for id, key := range plan.Keys {
		if path, ok := idToPath[id]; ok {
			keys[path] = key
		}
	}
	for _, id := range plan.Skip {
		if path, ok := idToPath[id]; ok {
			skipPaths = append(skipPaths, path)
		}
	}
	sort.Strings(skipPaths)
	return keys, skipPaths
}

// recordLeafCachePasses writes PutPass for leaves that passed.
// When allPassed is true, every key is stored. Otherwise failed paths are skipped.
// Errors are ignored (best-effort). Keys/failed use bare tree-relative paths
// for single-tree (same identity scheme as prepareLeafCache returns).
func recordLeafCachePasses(keys map[string]string, failed map[string]bool, allPassed bool) {
	if len(keys) == 0 {
		return
	}
	storeRoot, err := leafcache.ResolveStoreRoot()
	if err != nil {
		return
	}
	store, err := leafcache.NewStore(storeRoot)
	if err != nil {
		return
	}
	leafcache.RecordPasses(store, keys, failed, allPassed)
}

func passedCases(total, failCount int) int {
	passed := total - failCount
	if passed < 0 {
		return 0
	}
	return passed
}

// displayGoArgs returns user-visible go command arguments, omitting internal flags.
func displayGoArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		switch a {
		case "-mod=mod", "-json":
			continue
		}
		out = append(out, a)
	}
	return out
}
