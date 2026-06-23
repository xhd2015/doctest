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

	var cases []core.TreeCase
	var err error

	cases, err = core.DiscoverTreeCases(dir)
	if err != nil {
		return TestRunStats{}, err
	}
	if opts.SubDir != "" {
		cases = core.FilterBySubDir(cases, dir, opts.SubDir)
	}
	if opts.ChangedOnly {
		gitRoot, changedFiles, err := core.ChangedGitFiles(dir)
		if err != nil {
			return TestRunStats{}, err
		}
		cases = core.FilterByChangedFiles(cases, dir, gitRoot, changedFiles)
		if len(cases) == 0 {
			fmt.Fprintln(w, core.NoTestsChangedMessage)
			return TestRunStats{NoTestsChanged: true}, nil
		}
	}
	if len(cases) == 0 {
		return TestRunStats{}, fmt.Errorf("%s: no runnable test cases found", dir)
	}

	stats := TestRunStats{Total: len(cases)}

	ctx, err := newGenerateContext(dir, opts, cases, w, false, opts.Verbose)
	if err != nil {
		return TestRunStats{}, err
	}
	ctx.installInterruptCleanup()
	defer ctx.Close()

	if opts.Verbose {
		ctx.announceRoots()
		fmt.Fprintf(w, "doctest: %s\n\n", pathfmt.Short(dir))
		if _, err := core.DiscoverTreeCasesVerbose(dir, w); err != nil {
			return TestRunStats{}, err
		}
		fmt.Fprintf(w, "─── %d test cases\n\n", len(cases))
	} else {
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

	testArgs := []string{"test", "-mod=mod"}
	if opts.Verbose {
		testArgs = append(testArgs, "-v")
	}
	if NeedsBuildVCSFlag(runDir) {
		testArgs = append(testArgs, "-buildvcs=false")
	}
	testArgs = append(testArgs, packageArgs...)
	if opts.Count > 0 {
		testArgs = append(testArgs, fmt.Sprintf("-count=%d", opts.Count))
	}
	if opts.Timeout > 0 {
		testArgs = append(testArgs, fmt.Sprintf("-timeout=%s", opts.Timeout))
	}

	execArgs := append([]string(nil), testArgs...)
	if !opts.Verbose {
		execArgs = append(execArgs, "-json")
	}

	displayArgs := displayGoArgs(testArgs)
	if opts.Verbose {
		fmt.Fprintf(w, "cd %s && go %s\n\n", pathfmt.Short(runDir), strings.Join(displayArgs, " "))
	} else {
		fmt.Fprintf(w, "cd %s && go %s\n", pathfmt.Short(runDir), strings.Join(displayArgs, " "))
	}

	sessionID := core.DoctestSessionIDForRun()
	goTestCmd := exec.Command("go", execArgs...)
	goTestCmd.Dir = runDir
	goTestCmd.Env = append(os.Environ(), core.DoctestSessionIDEnv+"="+sessionID)

	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	if opts.Verbose {
		start := time.Now()
		out, err := goTestCmd.CombinedOutput()
		elapsed := time.Since(start)
		stdout.Write(out)
		stats.Passed = passedCases(stats.Total, countFailuresFromGoTestOutput(out))
		if !opts.SuppressResultSummary {
			stats.Elapsed = elapsed
			PrintResultSummary(opts, stats)
		}
		if err != nil {
			return stats, fmt.Errorf("go test failed: %v", err)
		}
	} else {
		stdoutPipe, err := goTestCmd.StdoutPipe()
		if err != nil {
			return TestRunStats{}, fmt.Errorf("stdout pipe: %w", err)
		}
		stderrPipe, err := goTestCmd.StderrPipe()
		if err != nil {
			return TestRunStats{}, fmt.Errorf("stderr pipe: %w", err)
		}

		if err := goTestCmd.Start(); err != nil {
			return TestRunStats{}, fmt.Errorf("go test start: %w", err)
		}
		start := time.Now()

		passCount := 0
		failCount := 0
		cachedCount := 0
		var failLines []string
		var detailLines []string
		style := newColorStyle(opts.Color, stdout)
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
					detailLines = append(detailLines, buf...)
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
						failLines = append(failLines, trimmed)
					}
					if trimmed != "" && trimmed != "PASS" {
						if ev.Test != "" {
							key := testKey(ev.Package, ev.Test)
							line := strings.TrimRight(ev.Output, "\n")
							if strings.Contains(ev.Output, "--- FAIL:") {
								flushTestOutput(key)
								detailLines = append(detailLines, line)
							} else if failedTests[key] {
								detailLines = append(detailLines, line)
							} else {
								testOutputs[key] = append(testOutputs[key], line)
							}
						} else if !strings.HasPrefix(trimmed, "ok ") && !strings.HasPrefix(trimmed, "ok\t") &&
							!strings.HasPrefix(trimmed, "FAIL\t") && !strings.HasPrefix(trimmed, "FAIL ") {
							detailLines = append(detailLines, trimmed)
						}
					}
				case "pass":
					if ev.Test != "" {
						delete(testOutputs, testKey(ev.Package, ev.Test))
						continue
					}
					passCount++
					if packageCached[ev.Package] {
						cachedCount++
					}
					stdout.Write([]byte("."))
				case "fail":
					if ev.Test != "" {
						key := testKey(ev.Package, ev.Test)
						flushTestOutput(key)
						continue
					}
					failCount++
					fmt.Fprint(stdout, style.red("."))
				}
			}
		}()

		stderrData, _ := io.ReadAll(stderrPipe)
		stdoutWg.Wait()

		err = goTestCmd.Wait()
		elapsed := time.Since(start)
		stats.Passed = passedCases(stats.Total, failCount)

		fmt.Fprintln(stdout, formatSummary(style, passCount+failCount, passCount, failCount, cachedCount, elapsed))

		if !opts.SuppressResultSummary {
			stats.Elapsed = elapsed
			PrintResultSummary(opts, stats)
		}

		if len(failLines) > 0 {
			for _, line := range failLines {
				fmt.Fprintln(stdout, line)
			}
		}
		for _, line := range detailLines {
			fmt.Fprintln(stdout, line)
		}

		if len(stderrData) > 0 {
			stdout.Write(stderrData)
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