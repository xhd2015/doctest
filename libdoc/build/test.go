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

	var cases []core.TreeCase
	var err error

	absRoot, _ := filepath.Abs(dir)

	mappingGenRoot := opts.GenDir
	if mappingGenRoot == "" {
		var err error
		mappingGenRoot, _, err = core.CacheMappingGenRoot(absRoot)
		if err != nil {
			return err
		}
	}
	if err := os.MkdirAll(mappingGenRoot, 0755); err != nil {
		return err
	}

	if opts.Stderr != nil {
		fmt.Fprintf(opts.Stderr, "→ %s\n\n", mappingGenRoot)
	} else {
		fmt.Fprintf(os.Stderr, "→ %s\n\n", mappingGenRoot)
	}

	if opts.Verbose {
		fmt.Fprintf(w, "doctest: %s\n\n", dir)
		cases, err = core.DiscoverTreeCasesVerbose(dir, w)
		if err != nil {
			return err
		}
		if opts.SubDir != "" {
			cases = core.FilterBySubDir(cases, dir, opts.SubDir)
		}
		fmt.Fprintf(w, "─── %d test cases\n\n", len(cases))
	} else {
		cases, err = core.DiscoverTreeCases(dir)
		if err != nil {
			return err
		}
		if opts.SubDir != "" {
			cases = core.FilterBySubDir(cases, dir, opts.SubDir)
		}
		fmt.Fprintf(w, "doctest: %s\n", dir)
		fmt.Fprintf(w, "─── %d test cases\n", len(cases))
	}

	if len(cases) == 0 {
		return fmt.Errorf("%s: no runnable test cases found", dir)
	}

	_, modPath, hasMod := core.FindModuleRoot(absRoot)
	absModRoot, _ := core.MappingGenRoot(absRoot)

	if err := core.WriteGoMod(mappingGenRoot, absModRoot, modPath, hasMod); err != nil {
		return err
	}
	if opts.Verbose {
		fmt.Fprintf(w, "→ %s\n", filepath.Join(mappingGenRoot, "go.mod"))
	}

	pkgName := "testcase"
	srcDir, origPkg, hasPkgUnderTest := core.ResolvePkgUnderTest(absRoot)
	if hasPkgUnderTest {
		pkgName = origPkg + "_tc"
	}

	for _, tc := range cases {
		absLeafDir := filepath.Join(absRoot, tc.Path)
		leafDir, err := core.GenDirForLeaf(mappingGenRoot, absModRoot, absLeafDir)
		if err != nil {
			return fmt.Errorf("gen dir for leaf %s: %w", tc.Path, err)
		}

		if hasPkgUnderTest {
			if _, err := core.CopySourceFiles(leafDir, srcDir, origPkg); err != nil {
				return fmt.Errorf("copy source files to %s: %w", leafDir, err)
			}
		}

		testPath, err := core.WriteGeneratedCase(leafDir, tc, false, pkgName, absRoot)
		if err != nil {
			return err
		}
		if opts.Verbose {
			fmt.Fprintf(w, "→ %s\n", testPath)
		}
	}

	if hasMod {
		if err := core.CondTidyGoMod(mappingGenRoot); err != nil {
			return err
		}
	}

	runDir := mappingGenRoot
	isSingleLeaf := false
	if opts.SubDir != "" {
		subDirAbs := opts.SubDir
		if !filepath.IsAbs(subDirAbs) {
			subDirAbs = filepath.Join(absRoot, subDirAbs)
		}
		if _, err := os.Stat(filepath.Join(subDirAbs, "ASSERT.md")); err == nil {
			isSingleLeaf = true
		}
		relSubDir, err := filepath.Rel(absModRoot, subDirAbs)
		if err == nil && relSubDir != "." {
			runDir = filepath.Join(mappingGenRoot, relSubDir)
		}
	}
	if !isSingleLeaf && len(cases) == 1 && cases[0].Path != "" {
		leafDir, _ := core.GenDirForLeaf(mappingGenRoot, absModRoot, filepath.Join(absRoot, cases[0].Path))
		runDir = leafDir
		isSingleLeaf = true
	}

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
				} else if strings.HasPrefix(line, "FAIL\t") || strings.HasPrefix(line, "FAIL ") {
					runCount++
					failCount++
				}
			}
			if scanner.Err() != nil {
				fmt.Fprintf(os.Stderr, "read go test stdout: %v\n", scanner.Err())
			}
		}()

		stderrData, _ := io.ReadAll(stderrPipe)
		stdoutWg.Wait()

		err = goTestCmd.Wait()

		for i := 0; i < runCount; i++ {
			fmt.Print(".")
		}
		fmt.Printf("  (%d Run, %d Pass, %d Fail)\n", runCount, passCount, failCount)

		if len(stderrData) > 0 {
			os.Stdout.Write(stderrData)
		}

		if err != nil {
			return fmt.Errorf("go test: %w", err)
		}
	}
	return nil
}
