package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Test(dir string, opts core.Options) error {
	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}

	var cases []core.TreeCase
	var err error

	tmp := opts.GenDir
	removeTmp := false
	if tmp == "" {
		tmp, err = os.MkdirTemp("", "doctest-test-*")
		if err != nil {
			return err
		}
		removeTmp = opts.RemoveTemp
	} else if err := os.MkdirAll(tmp, 0755); err != nil {
		return err
	}

	fmt.Fprintf(w, "→ %s\n\n", tmp)

	if removeTmp {
		defer os.RemoveAll(tmp)
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

	modRoot, modPath, hasMod := core.FindModuleRoot(dir)
	if err := core.WriteGoMod(tmp, modRoot, modPath, hasMod); err != nil {
		return err
	}
	if opts.Verbose {
		fmt.Fprintf(w, "→ %s\n", filepath.Join(tmp, "go.mod"))
	}

	pkgName := "testcase"
	if srcDir, origPkg, ok := core.ResolvePkgUnderTest(dir); ok {
		newPkg, err := core.CopySourceFiles(tmp, srcDir, origPkg)
		if err != nil {
			return fmt.Errorf("copy source files: %w", err)
		}
		pkgName = newPkg
		if opts.Verbose {
			fmt.Fprintf(w, "→ %s (copied from %s, package %s)\n", srcDir, srcDir, newPkg)
		}
	}

	absRoot, _ := filepath.Abs(dir)
	_, err = core.WriteGeneratedCases(tmp, cases, false, nil, pkgName, absRoot)
	if err != nil {
		return err
	}

	testBinPath := filepath.Join(tmp, "test.bin")

	fmt.Fprintf(w, "cd %s && go test -c -o test.bin . && cd %s && %s/test.bin\n\n", tmp, absRoot, tmp)

	buildArgs := []string{"test", "-c", "-mod=mod"}
	if NeedsBuildVCSFlag(tmp) {
		buildArgs = append(buildArgs, "-buildvcs=false")
	}
	buildArgs = append(buildArgs, "-o", testBinPath, ".")
	goTestBuild := exec.Command("go", buildArgs...)
	goTestBuild.Dir = tmp
	if opts.Verbose {
		goTestBuild.Stdout = w
		goTestBuild.Stderr = w
		if err := goTestBuild.Run(); err != nil {
			return fmt.Errorf("go test -c failed: %v", err)
		}
	} else {
		if out, err := goTestBuild.CombinedOutput(); err != nil {
			return fmt.Errorf("go test -c failed: %v\n%s", err, string(out))
		}
	}

	args := []string{}
	if opts.Count > 0 {
		args = append(args, fmt.Sprintf("-test.count=%d", opts.Count))
	}
	if opts.Verbose {
		args = append(args, "-test.v")
	}
	runCmd := exec.Command(testBinPath, args...)
	runCmd.Dir = absRoot
	if opts.Verbose {
		runCmd.Stdout = w
		runCmd.Stderr = w
		if err := runCmd.Run(); err != nil {
			return fmt.Errorf("test binary failed: %v", err)
		}
		return nil
	}
	if out, err := runCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("test binary failed: %v\n%s", err, string(out))
	}
	return nil
}
