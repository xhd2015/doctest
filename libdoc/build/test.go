package build

import (
	"bufio"
	"bytes"
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
				PrintSkippedSummary(skipped, opts.Verbose)
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
	if opts.Timeout > 0 {
		flagArgs = append(flagArgs, fmt.Sprintf("-timeout=%s", opts.Timeout))
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

	displayArgs := displayGoArgs(append(append([]string(nil), flagArgs...), packageArgs...))
	if opts.Verbose {
		fmt.Fprintf(w, "cd %s && go %s\n\n", pathfmt.Short(runDir), strings.Join(displayArgs, " "))
	} else {
		fmt.Fprintf(w, "cd %s && go %s\n", pathfmt.Short(runDir), strings.Join(displayArgs, " "))
	}

	sessionID := core.DoctestSessionIDForRun()
	goCache := opts.GoCache
	// Prefer explicit nest sink on opts; else inherit process (suite child).
	if opts.MetricsNestSink == "" {
		if v := os.Getenv(envMetricsNestSink); v != "" {
			opts.MetricsNestSink = v
		}
	}

	// Programmatic leaf-cache: warm GetPass hits skip leaf bodies via suite env.
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
	tGo := time.Now()
	var goTestElapsed time.Duration
	if opts.Verbose {
		var out []byte
		var err error
		if singlePkgInvocations {
			for _, pkg := range packageArgs {
				execArgs := append(append([]string(nil), flagArgs...), pkg)
				goTestCmd := exec.Command("go", execArgs...)
				goTestCmd.Dir = runDir
				goTestCmd.Env = goTestEnvFull(sessionID, goCache, opts.MetricsNestSink, "", leafSkipEnv)
				pkgOut, pkgErr := goTestCmd.CombinedOutput()
				out = append(out, pkgOut...)
				if pkgErr != nil && err == nil {
					err = pkgErr
				}
			}
		} else {
			execArgs := append(append([]string(nil), flagArgs...), packageArgs...)
			goTestCmd := exec.Command("go", execArgs...)
			goTestCmd.Dir = runDir
			goTestCmd.Env = goTestEnvFull(sessionID, goCache, opts.MetricsNestSink, "", leafSkipEnv)
			out, err = goTestCmd.CombinedOutput()
		}
		goTestElapsed = time.Since(tGo)
		track("go_test", tGo)
		stdout.Write(out)
		stats.Passed = passedCases(stats.Total, countFailuresFromGoTestOutput(out))
		// Verbose path: no package Elapsed; single-leaf gets full go_test wall.
		if len(cases) == 1 {
			stats.LeafTimings = []LeafTiming{{Path: cases[0].Path, ElapsedNs: goTestElapsed.Nanoseconds()}}
		}
		// Best-effort PutPass: verbose mode lacks per-leaf fail map — store only on full success.
		if err == nil {
			recordLeafCachePasses(leafKeys, nil, true)
		}
		if !opts.SuppressResultSummary {
			stats.Elapsed = goTestElapsed
			PrintSkippedSummary(stats.Skipped, opts.Verbose)
			PrintResultSummary(opts, stats)
		}
		if err != nil {
			stats.Phases = phases
			return stats, fmt.Errorf("go test failed: %v", err)
		}
	} else {
		style := newColorStyle(opts.Color, stdout)
		var result goTestJSONResult
		var err error
		if singlePkgInvocations {
			result, err = runGoTestJSONPerPackage(runDir, flagArgs, packageArgs, sessionID, goCache, opts.MetricsNestSink, leafSkipEnv, stdout, style)
		} else {
			result, err = runGoTestJSONShards(runDir, flagArgs, packageArgs, sessionID, goCache, opts.MetricsNestSink, leafSkipEnv, stdout, style)
		}
		goTestElapsed = time.Since(tGo)
		track("go_test", tGo)
		stats.Passed = passedCases(stats.Total, result.failCount)
		if ctx.unifiedMode {
			stats.LeafTimings = leafTimingsFromSubtests(cases, result, goTestElapsed)
		} else {
			stats.LeafTimings = leafTimingsFromPackages(cases, packageArgs, isSingleLeaf, result, goTestElapsed)
		}

		// Programmatic leaf-cache skips count toward summary Cached when the
		// go test package itself was not (cached).
		if !result.anyPackageCached() {
			result.cachedCount += len(skipPaths)
		}
		recordLeafCachePasses(leafKeys, result.suiteLeafFailed, err == nil && result.failCount == 0)

		fmt.Fprintln(stdout, formatSummary(style, result.passCount+result.failCount, result.passCount, result.failCount, result.cachedCount, goTestElapsed))

		if !opts.SuppressResultSummary {
			stats.Elapsed = goTestElapsed
			PrintSkippedSummary(stats.Skipped, opts.Verbose)
			PrintResultSummary(opts, stats)
		}

		for _, line := range result.failLines {
			fmt.Fprintln(stdout, line)
		}
		for _, line := range result.detailLines {
			fmt.Fprintln(stdout, line)
		}
		if len(result.stderrData) > 0 {
			stdout.Write(result.stderrData)
		}

		if err != nil {
			stats.Phases = phases
			return stats, fmt.Errorf("go test: %w", err)
		}
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
	passCount     int
	failCount     int
	cachedCount   int
	failLines     []string
	detailLines   []string
	stderrData    []byte
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
	dst.cachedCount += src.cachedCount
	dst.failLines = append(dst.failLines, src.failLines...)
	dst.detailLines = append(dst.detailLines, src.detailLines...)
	dst.stderrData = append(dst.stderrData, src.stderrData...)
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
// goTestEnvWithOpts over mid-suite os.Setenv.
func goTestEnv(sessionID, goCache string) []string {
	return goTestEnvFull(sessionID, goCache, "", "", "")
}

func goTestEnvFull(sessionID, goCache, nestSink, goWork, leafSkipPaths string) []string {
	env := append(os.Environ(), core.DoctestSessionIDEnv+"="+sessionID)
	if goCache != "" {
		env = append(env, "GOCACHE="+goCache)
	}
	if nestSink != "" {
		env = append(env, envMetricsNestSink+"="+nestSink)
	} else if v := os.Getenv(envMetricsNestSink); v != "" {
		// Inherit outer nest sink for nested go test (read-only).
		env = append(env, envMetricsNestSink+"="+v)
	}
	if goWork != "" {
		// Multi-module workspace hub (go.work). Empty means use process default / off.
		env = append(env, "GOWORK="+goWork)
	}
	if leafSkipPaths != "" {
		env = append(env, leafcache.EnvSkipPaths+"="+leafSkipPaths)
	}
	return env
}

func goTestEnvFromOpts(sessionID string, opts core.Options) []string {
	return goTestEnvFull(sessionID, opts.GoCache, opts.MetricsNestSink, "", "")
}

// runGoTestJSONPerPackage runs one go test process per package (serial) so
// profile flags that go rejects with multi-package lists still work.
func runGoTestJSONPerPackage(runDir string, flagArgs, packageArgs []string, sessionID, goCache, nestSink, leafSkipPaths string, stdout io.Writer, style colorStyle) (goTestJSONResult, error) {
	var merged goTestJSONResult
	var firstErr error
	for _, pkg := range packageArgs {
		args := append(append([]string(nil), flagArgs...), pkg)
		res, err := runGoTestJSONOnce(runDir, args, sessionID, goCache, nestSink, "", leafSkipPaths, stdout, style)
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

func runGoTestJSONShards(runDir string, flagArgs, packageArgs []string, sessionID, goCache, nestSink, leafSkipPaths string, stdout io.Writer, style colorStyle) (goTestJSONResult, error) {
	// Single go test process per tree. Package sharding multiplies nested
	// self-test fan-out and has raced go.mod; wall cut is tree concurrency +
	// heavy/light scheduling in path_resolve.
	workers := 1
	shards := packageTestShards(packageArgs, workers)
	if len(shards) <= 1 {
		args := append(append([]string(nil), flagArgs...), packageArgs...)
		return runGoTestJSONOnce(runDir, args, sessionID, goCache, nestSink, "", leafSkipPaths, stdout, style)
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
			res, err := runGoTestJSONOnce(runDir, args, sessionID, goCache, nestSink, "", leafSkipPaths, &lockedWriter{w: stdout, mu: &mu}, style)
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

func runGoTestJSONOnce(runDir string, testArgs []string, sessionID, goCache, nestSink, goWork, leafSkipPaths string, stdout io.Writer, style colorStyle) (goTestJSONResult, error) {
	goTestSlots <- struct{}{}
	defer func() { <-goTestSlots }()

	unlockMod := lockGoTestModule(runDir)
	defer unlockMod()

	execArgs := append(append([]string(nil), testArgs...), "-json")
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
		// package-level dots and use leaf pass/fail/cached for the summary so
		// one suite binary still reports one result per leaf (not 1 package).
		suiteLeafPass := 0
		suiteLeafFail := 0
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
		countSuiteLeaf := func(pkg, test string, failed bool) {
			key := testKey(pkg, test)
			if suiteLeafSeen[key] {
				return
			}
			suiteLeafSeen[key] = true
			if failed {
				suiteLeafFail++
				if rel := suiteLeafRelPath(test); rel != "" {
					res.suiteLeafFailed[rel] = true
				}
			} else {
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
				trimmed := strings.TrimSpace(ev.Output)
				if ev.Test == "" && (strings.HasPrefix(trimmed, "ok ") || strings.HasPrefix(trimmed, "ok\t")) {
					if strings.Contains(ev.Output, "(cached)") {
						packageCached[ev.Package] = true
					}
				}
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
			case "pass":
				if ev.Test != "" {
					recordTest(ev.Test, ev.Elapsed)
					delete(testOutputs, testKey(ev.Package, ev.Test))
					if isCountableSuiteLeaf(ev.Package, ev.Test) {
						before := suiteLeafPass + suiteLeafFail
						countSuiteLeaf(ev.Package, ev.Test, false)
						if suiteLeafPass+suiteLeafFail > before {
							stdout.Write([]byte("."))
						}
					}
					continue
				}
				res.passCount++
				if packageCached[ev.Package] {
					res.cachedCount++
				}
				recordPkg(ev.Package, ev.Elapsed, packageCached[ev.Package])
				if suiteLeafPass+suiteLeafFail == 0 {
					stdout.Write([]byte("."))
				}
			case "fail":
				if ev.Test != "" {
					recordTest(ev.Test, ev.Elapsed)
					key := testKey(ev.Package, ev.Test)
					flushTestOutput(key)
					if isCountableSuiteLeaf(ev.Package, ev.Test) {
						before := suiteLeafPass + suiteLeafFail
						countSuiteLeaf(ev.Package, ev.Test, true)
						if suiteLeafPass+suiteLeafFail > before {
							fmt.Fprint(stdout, style.red("."))
						}
					}
					continue
				}
				res.failCount++
				recordPkg(ev.Package, ev.Elapsed, false)
				if suiteLeafPass+suiteLeafFail == 0 {
					fmt.Fprint(stdout, style.red("."))
				}
			}
		}
		// Prefer leaf-level counts for unified suite so multi-leaf trees report
		// N Run / N Pass instead of a single package result.
		if suiteLeafPass+suiteLeafFail > 0 {
			res.passCount = suiteLeafPass
			res.failCount = suiteLeafFail
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
	if err != nil {
		return res, err
	}
	return res, nil
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

// prepareLeafCache computes leaf keys and warm skip paths for this tree run.
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
	for _, tc := range cases {
		in, err := leafcache.KeyForLeaf(treeRoot, tc.Path, goVer)
		if err != nil {
			continue
		}
		key, err := leafcache.ComputeLeafKey(in)
		if err != nil {
			continue
		}
		keys[tc.Path] = key
		if !enabled {
			continue
		}
		hit, err := store.GetPass(key)
		if err != nil || !hit {
			continue
		}
		skipPaths = append(skipPaths, tc.Path)
	}
	sort.Strings(skipPaths)
	return keys, skipPaths
}

// recordLeafCachePasses writes PutPass for leaves that passed.
// When allPassed is true, every key is stored. Otherwise failed paths are skipped.
// Errors are ignored (best-effort).
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
	for path, key := range keys {
		if !allPassed {
			if failed != nil && failed[path] {
				continue
			}
			// Partial fail without per-leaf map: only store when allPassed or known non-fail.
			// When failed map is nil and !allPassed, skip all to avoid storing fails.
			if failed == nil {
				continue
			}
		}
		_ = store.PutPass(key)
	}
}

func countFailuresFromGoTestOutput(out []byte) int {
	failures := 0
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "FAIL\t") || strings.HasPrefix(line, "FAIL ") {
			failures++
		}
	}
	return failures
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