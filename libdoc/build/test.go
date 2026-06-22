package build

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/pathfmt"
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
		fmt.Fprintf(w, "doctest: %s\n\n", pathfmt.DisplayPath(dir))
		if _, err := core.DiscoverTreeCasesVerbose(dir, w); err != nil {
			return TestRunStats{}, err
		}
		fmt.Fprintf(w, "─── %d test cases\n\n", len(cases))
	} else {
		fmt.Fprintf(w, "doctest: %s (%d tests)\n", pathfmt.DisplayPath(dir), len(cases))
	}

	if err := ctx.writeCases(cases, false); err != nil {
		return TestRunStats{}, err
	}

	runDir, isSingleLeaf := ctx.runDir(absRoot, opts, cases)

	testArgs := []string{"test", "-mod=mod"}
	if opts.Verbose {
		testArgs = append(testArgs, "-v")
	}
	if NeedsBuildVCSFlag(runDir) {
		testArgs = append(testArgs, "-buildvcs=false")
	}
	if isSingleLeaf {
		testArgs = append(testArgs, ".")
	} else {
		testArgs = append(testArgs, "./...")
	}
	if opts.Count > 0 {
		testArgs = append(testArgs, fmt.Sprintf("-count=%d", opts.Count))
	}
	if opts.Timeout > 0 {
		testArgs = append(testArgs, fmt.Sprintf("-timeout=%s", opts.Timeout))
	}

	if opts.Verbose {
		fmt.Fprintf(w, "cd %s && go %s\n\n", pathfmt.DisplayPath(runDir), strings.Join(testArgs, " "))
	} else {
		fmt.Fprintf(w, "cd %s && go %s\n", pathfmt.DisplayPath(runDir), strings.Join(testArgs, " "))
	}

	sessionID := core.DoctestSessionIDForRun()
	goTestCmd := exec.Command("go", testArgs...)
	goTestCmd.Dir = runDir
	goTestCmd.Env = append(os.Environ(), core.DoctestSessionIDEnv+"="+sessionID)

	if opts.Verbose {
		start := time.Now()
		out, err := goTestCmd.CombinedOutput()
		elapsed := time.Since(start)
		os.Stdout.Write(out)
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
		style := newColorStyle(opts.Color, os.Stdout)
		var stdoutWg sync.WaitGroup
		stdoutWg.Add(1)
		go func() {
			defer stdoutWg.Done()
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "ok\t") {
					passCount++
					if strings.Contains(line, "(cached)") {
						cachedCount++
					}
					os.Stdout.Write([]byte("."))
				} else if strings.HasPrefix(line, "FAIL\t") || strings.HasPrefix(line, "FAIL ") {
					failCount++
					failLines = append(failLines, line)
					os.Stdout.WriteString(style.red("."))
				} else {
					detailLines = append(detailLines, line)
				}
			}
			if scanner.Err() != nil {
				fmt.Fprintf(os.Stderr, "read go test stdout: %v\n", scanner.Err())
			}
		}()

		stderrData, _ := io.ReadAll(stderrPipe)
		stdoutWg.Wait()

		err = goTestCmd.Wait()
		elapsed := time.Since(start)
		stats.Passed = passedCases(stats.Total, failCount)

		fmt.Println(formatSummary(style, passCount+failCount, passCount, failCount, cachedCount, elapsed))

		if !opts.SuppressResultSummary {
			stats.Elapsed = elapsed
			PrintResultSummary(opts, stats)
		}

		if len(failLines) > 0 {
			for _, line := range failLines {
				fmt.Println(line)
			}
		}
		for _, line := range detailLines {
			fmt.Println(line)
		}

		if len(stderrData) > 0 {
			os.Stdout.Write(stderrData)
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