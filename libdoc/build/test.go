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
	"strings"
	"sync"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
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
	tDiscover := time.Now()
	allCases, err := core.DiscoverTreeCases(dir)
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
	track("discover", tDiscover)
	if len(cases) == 0 {
		if len(skipped) > 0 {
			stats := TestRunStats{Skipped: skipped, Phases: phases}
			if !opts.SuppressResultSummary {
				PrintSkippedSummary(skipped)
			}
			return stats, nil
		}
		return TestRunStats{Phases: phases}, fmt.Errorf("%s: no runnable test cases found", dir)
	}

	stats := TestRunStats{Total: len(cases), Skipped: skipped}

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

	// Force ref when unified is set (also done at parse time for CLI).
	if opts.ExperimentUnifiedPackagePerDoctestTree {
		opts.ExperimentRefInsteadOfInline = true
		fmt.Fprintln(w, "doctest: experiment: unified-package-per-doctest-tree")
	}
	if opts.ExperimentRefInsteadOfInline {
		fmt.Fprintln(w, "doctest: experiment: ref-instead-of-inline")
	}

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

	runDir, isSingleLeaf := ctx.runDir(absRoot, opts, cases)

	var packageArgs []string
	if ctx.unifiedMode {
		// Always suite-only packaging (even for a single leaf).
		runDir = ctx.scopedMultiRunDir(absRoot)
		var pkgErr error
		packageArgs, pkgErr = ctx.packageArgsForCases(runDir, absRoot, cases)
		if pkgErr != nil {
			track("generate", tGen)
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
			track("generate", tGen)
			stats.Phases = phases
			return stats, pkgErr
		}
	}
	track("generate", tGen)

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
				goTestCmd.Env = goTestEnv(sessionID, goCache)
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
			goTestCmd.Env = goTestEnv(sessionID, goCache)
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
		if !opts.SuppressResultSummary {
			stats.Elapsed = goTestElapsed
			PrintSkippedSummary(stats.Skipped)
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
			result, err = runGoTestJSONPerPackage(runDir, flagArgs, packageArgs, sessionID, goCache, stdout, style)
		} else {
			result, err = runGoTestJSONShards(runDir, flagArgs, packageArgs, sessionID, goCache, stdout, style)
		}
		goTestElapsed = time.Since(tGo)
		track("go_test", tGo)
		stats.Passed = passedCases(stats.Total, result.failCount)
		if ctx.unifiedMode {
			stats.LeafTimings = leafTimingsFromSubtests(cases, result, goTestElapsed)
		} else {
			stats.LeafTimings = leafTimingsFromPackages(cases, packageArgs, isSingleLeaf, result, goTestElapsed)
		}

		fmt.Fprintln(stdout, formatSummary(style, result.passCount+result.failCount, result.passCount, result.failCount, result.cachedCount, goTestElapsed))

		if !opts.SuppressResultSummary {
			stats.Elapsed = goTestElapsed
			PrintSkippedSummary(stats.Skipped)
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
// unified suite (t.Run(leafPath, …) → TestDoctestSuite/<leafPath>).
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

// goTestEnv builds the environment for a child `go test` process, including
// DOCTEST_SESSION_ID and an optional isolated GOCACHE (cold-cache mode).
func goTestEnv(sessionID, goCache string) []string {
	env := append(os.Environ(), core.DoctestSessionIDEnv+"="+sessionID)
	if goCache != "" {
		env = append(env, "GOCACHE="+goCache)
	}
	return env
}

// runGoTestJSONPerPackage runs one go test process per package (serial) so
// profile flags that go rejects with multi-package lists still work.
func runGoTestJSONPerPackage(runDir string, flagArgs, packageArgs []string, sessionID, goCache string, stdout io.Writer, style colorStyle) (goTestJSONResult, error) {
	var merged goTestJSONResult
	var firstErr error
	for _, pkg := range packageArgs {
		args := append(append([]string(nil), flagArgs...), pkg)
		res, err := runGoTestJSONOnce(runDir, args, sessionID, goCache, stdout, style)
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

func runGoTestJSONShards(runDir string, flagArgs, packageArgs []string, sessionID, goCache string, stdout io.Writer, style colorStyle) (goTestJSONResult, error) {
	// Single go test process per tree. Package sharding multiplies nested
	// self-test fan-out and has raced go.mod; wall cut is tree concurrency +
	// heavy/light scheduling in path_resolve.
	workers := 1
	shards := packageTestShards(packageArgs, workers)
	if len(shards) <= 1 {
		args := append(append([]string(nil), flagArgs...), packageArgs...)
		return runGoTestJSONOnce(runDir, args, sessionID, goCache, stdout, style)
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
			res, err := runGoTestJSONOnce(runDir, args, sessionID, goCache, &lockedWriter{w: stdout, mu: &mu}, style)
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

func runGoTestJSONOnce(runDir string, testArgs []string, sessionID, goCache string, stdout io.Writer, style colorStyle) (goTestJSONResult, error) {
	goTestSlots <- struct{}{}
	defer func() { <-goTestSlots }()

	unlockMod := lockGoTestModule(runDir)
	defer unlockMod()

	execArgs := append(append([]string(nil), testArgs...), "-json")
	goTestCmd := exec.Command("go", execArgs...)
	goTestCmd.Dir = runDir
	goTestCmd.Env = goTestEnv(sessionID, goCache)

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
		testKey := func(pkg, test string) string { return pkg + "\x00" + test }
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
					continue
				}
				res.passCount++
				if packageCached[ev.Package] {
					res.cachedCount++
				}
				recordPkg(ev.Package, ev.Elapsed, packageCached[ev.Package])
				stdout.Write([]byte("."))
			case "fail":
				if ev.Test != "" {
					recordTest(ev.Test, ev.Elapsed)
					key := testKey(ev.Package, ev.Test)
					flushTestOutput(key)
					continue
				}
				res.failCount++
				recordPkg(ev.Package, ev.Elapsed, false)
				fmt.Fprint(stdout, style.red("."))
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