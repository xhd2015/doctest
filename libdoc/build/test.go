package build

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Test(dir string, opts core.Options) error {
	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}

	absRoot, _ := filepath.Abs(dir)

	var cases []core.TreeCase
	var err error

	cases, err = core.DiscoverTreeCases(dir)
	if err != nil {
		return err
	}
	if opts.SubDir != "" {
		cases = core.FilterBySubDir(cases, dir, opts.SubDir)
	}
	if len(cases) == 0 {
		return fmt.Errorf("%s: no runnable test cases found", dir)
	}

	ctx, err := newGenerateContext(dir, opts, cases, w, false, opts.Verbose)
	if err != nil {
		return err
	}
	defer ctx.Close()

	ctx.announceRoots()

	if opts.Verbose {
		fmt.Fprintf(w, "doctest: %s\n\n", dir)
		if _, err := core.DiscoverTreeCasesVerbose(dir, w); err != nil {
			return err
		}
		fmt.Fprintf(w, "─── %d test cases\n\n", len(cases))
	} else {
		fmt.Fprintf(w, "doctest: %s\n", dir)
		fmt.Fprintf(w, "─── %d test cases\n", len(cases))
	}

	if err := ctx.writeCases(cases, false); err != nil {
		return err
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

	fmt.Fprintf(w, "cd %s && go %s\n\n", runDir, strings.Join(testArgs, " "))

	goTestCmd := exec.Command("go", testArgs...)
	goTestCmd.Dir = runDir

	if opts.Verbose {
		out, err := goTestCmd.CombinedOutput()
		os.Stdout.Write(out)
		if err != nil {
			return fmt.Errorf("go test failed: %v", err)
		}
	} else {
		stdoutPipe, err := goTestCmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("stdout pipe: %w", err)
		}
		stderrPipe, err := goTestCmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("stderr pipe: %w", err)
		}

		if err := goTestCmd.Start(); err != nil {
			return fmt.Errorf("go test start: %w", err)
		}

		runCount := 0
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
					runCount++
					passCount++
					if strings.Contains(line, "(cached)") {
						cachedCount++
					}
					os.Stdout.Write([]byte("."))
				} else if strings.HasPrefix(line, "FAIL\t") || strings.HasPrefix(line, "FAIL ") {
					runCount++
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

		fmt.Println(formatSummary(style, runCount, passCount, failCount, cachedCount))

		if len(failLines) > 0 {
			fmt.Println()
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
			return fmt.Errorf("go test: %w", err)
		}
	}

	if err := ctx.syncDump(); err != nil {
		return err
	}
	return nil
}