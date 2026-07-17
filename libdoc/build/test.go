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

	var allCases, cases []core.TreeCase
	var err error

	allCases, err = core.DiscoverTreeCases(dir)
	if err != nil {
		return TestRunStats{}, err
	}
	if opts.SubDir != "" {
		allCases = core.FilterBySubDir(allCases, dir, opts.SubDir)
	}

	var changedInfo core.ChangedRunInfo
	cases = allCases
	if opts.ChangedOnly {
		gitRoot, changedFiles, err := core.ChangedGitFiles(dir)
		if err != nil {
			return TestRunStats{}, err
		}
		changedInfo = core.ChangedRunInfoForTree(allCases, dir, gitRoot, changedFiles)
		cases = core.FilterByChangedFiles(allCases, dir, gitRoot, changedFiles)
		if core.ShouldAnnounceChangedRun(changedInfo, opts.Verbose) {
			fmt.Fprintln(w, core.FormatDoctestAnnouncement(pathfmt.Short(dir), changedInfo, true, 0))
		}
		if len(cases) == 0 {
			return TestRunStats{NoTestsChanged: true}, nil
		}
	}

	cases, skipped := core.FilterCasesByLabel(cases, opts)
	for i := range skipped {
		skipped[i].DisplayPath = SkippedDisplayPath(dir, skipped[i].Path)
	}
	if len(cases) == 0 {
		if len(skipped) > 0 {
			stats := TestRunStats{Skipped: skipped}
			if !opts.SuppressResultSummary {
				PrintSkippedSummary(skipped)
			}
			return stats, nil
		}
		return TestRunStats{}, fmt.Errorf("%s: no runnable test cases found", dir)
	}

	stats := TestRunStats{Total: len(cases), Skipped: skipped}

	ctx, err := newGenerateContext(dir, opts, cases, w, false, opts.Verbose)
	if err != nil {
		return TestRunStats{}, err
	}
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
			return TestRunStats{}, err
		}
		fmt.Fprintf(w, "─── %d test cases\n\n", len(cases))
	} else if !opts.ChangedOnly {
		fmt.Fprintf(w, "doctest: %s (%d tests)\n", pathfmt.Short(dir), len(cases))
	}

	if err := ctx.writeCases(cases, false); err != nil {
		return TestRunStats{}, err
	}

	runDir, isSingleLeaf := ctx.runDir(absRoot, opts, cases)

	var packageArgs []string
	if isSingleLeaf {
		packageArgs = []string{"."}
	} else {
		runDir = ctx.scopedMultiRunDir(absRoot)
		var pkgErr error
		packageArgs, pkgErr = ctx.packageArgsForCases(runDir, absRoot, cases)
		if pkgErr != nil {
			return TestRunStats{}, pkgErr
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

	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	// go test rejects most profile/coverprofile/trace flags when multiple packages
	// are listed on one invocation. Run packages one-at-a-time when those flags are set.
	singlePkgInvocations := profileFlagsNeedSinglePackage(opts) && len(packageArgs) > 1

	if opts.Verbose {
		start := time.Now()
		var out []byte
		var err error
		if singlePkgInvocations {
			for _, pkg := range packageArgs {
				execArgs := append(append([]string(nil), flagArgs...), pkg)
				goTestCmd := exec.Command("go", execArgs...)
				goTestCmd.Dir = runDir
				goTestCmd.Env = append(os.Environ(), core.DoctestSessionIDEnv+"="+sessionID)
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
			goTestCmd.Env = append(os.Environ(), core.DoctestSessionIDEnv+"="+sessionID)
			out, err = goTestCmd.CombinedOutput()
		}
		elapsed := time.Since(start)
		stdout.Write(out)
		stats.Passed = passedCases(stats.Total, countFailuresFromGoTestOutput(out))
		if !opts.SuppressResultSummary {
			stats.Elapsed = elapsed
			PrintSkippedSummary(stats.Skipped)
			PrintResultSummary(opts, stats)
		}
		if err != nil {
			return stats, fmt.Errorf("go test failed: %v", err)
		}
	} else {
		start := time.Now()
		style := newColorStyle(opts.Color, stdout)
		var result goTestJSONResult
		var err error
		if singlePkgInvocations {
			result, err = runGoTestJSONPerPackage(runDir, flagArgs, packageArgs, sessionID, stdout, style)
		} else {
			result, err = runGoTestJSONShards(runDir, flagArgs, packageArgs, sessionID, stdout, style)
		}
		elapsed := time.Since(start)
		stats.Passed = passedCases(stats.Total, result.failCount)

		fmt.Fprintln(stdout, formatSummary(style, result.passCount+result.failCount, result.passCount, result.failCount, result.cachedCount, elapsed))

		if !opts.SuppressResultSummary {
			stats.Elapsed = elapsed
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
			return stats, fmt.Errorf("go test: %w", err)
		}
	}

	if err := ctx.syncDump(); err != nil {
		return stats, err
	}
	if !opts.Verbose {
		fmt.Fprintln(w)
	}
	return stats, nil
}

type goTestJSONResult struct {
	passCount   int
	failCount   int
	cachedCount int
	failLines   []string
	detailLines []string
	stderrData  []byte
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

// runGoTestJSONPerPackage runs one go test process per package (serial) so
// profile flags that go rejects with multi-package lists still work.
func runGoTestJSONPerPackage(runDir string, flagArgs, packageArgs []string, sessionID string, stdout io.Writer, style colorStyle) (goTestJSONResult, error) {
	var merged goTestJSONResult
	var firstErr error
	for _, pkg := range packageArgs {
		args := append(append([]string(nil), flagArgs...), pkg)
		res, err := runGoTestJSONOnce(runDir, args, sessionID, stdout, style)
		merged.passCount += res.passCount
		merged.failCount += res.failCount
		merged.cachedCount += res.cachedCount
		merged.failLines = append(merged.failLines, res.failLines...)
		merged.detailLines = append(merged.detailLines, res.detailLines...)
		merged.stderrData = append(merged.stderrData, res.stderrData...)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return merged, firstErr
}

// goTestSlots bounds concurrent `go test` processes in this doctest process.
// Nested self-test subprocesses have their own pool (separate process).
var goTestSlots = make(chan struct{}, 4)

func runGoTestJSONShards(runDir string, flagArgs, packageArgs []string, sessionID string, stdout io.Writer, style colorStyle) (goTestJSONResult, error) {
	// Single go test process per tree. Package sharding multiplies nested
	// self-test fan-out and has raced go.mod; wall cut is tree concurrency +
	// heavy/light scheduling in path_resolve.
	workers := 1
	shards := packageTestShards(packageArgs, workers)
	if len(shards) <= 1 {
		args := append(append([]string(nil), flagArgs...), packageArgs...)
		return runGoTestJSONOnce(runDir, args, sessionID, stdout, style)
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
			res, err := runGoTestJSONOnce(runDir, args, sessionID, &lockedWriter{w: stdout, mu: &mu}, style)
			mu.Lock()
			defer mu.Unlock()
			merged.passCount += res.passCount
			merged.failCount += res.failCount
			merged.cachedCount += res.cachedCount
			merged.failLines = append(merged.failLines, res.failLines...)
			merged.detailLines = append(merged.detailLines, res.detailLines...)
			merged.stderrData = append(merged.stderrData, res.stderrData...)
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

func runGoTestJSONOnce(runDir string, testArgs []string, sessionID string, stdout io.Writer, style colorStyle) (goTestJSONResult, error) {
	goTestSlots <- struct{}{}
	defer func() { <-goTestSlots }()

	execArgs := append(append([]string(nil), testArgs...), "-json")
	goTestCmd := exec.Command("go", execArgs...)
	goTestCmd.Dir = runDir
	goTestCmd.Env = append(os.Environ(), core.DoctestSessionIDEnv+"="+sessionID)

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
	var stdoutWg sync.WaitGroup
	stdoutWg.Add(1)
	go func() {
		defer stdoutWg.Done()
		type goTestEvent struct {
			Action  string `json:"Action"`
			Package string `json:"Package"`
			Test    string `json:"Test"`
			Output  string `json:"Output"`
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
					delete(testOutputs, testKey(ev.Package, ev.Test))
					continue
				}
				res.passCount++
				if packageCached[ev.Package] {
					res.cachedCount++
				}
				stdout.Write([]byte("."))
			case "fail":
				if ev.Test != "" {
					key := testKey(ev.Package, ev.Test)
					flushTestOutput(key)
					continue
				}
				res.failCount++
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