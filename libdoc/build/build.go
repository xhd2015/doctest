package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Build(dir string, opts core.Options) error {
	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}

	var cases []core.TreeCase
	var err error

	tmp := opts.GenDir
	removeTmp := false
	if tmp == "" {
		tmp, err = os.MkdirTemp("", "doctest-build-*")
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
	_, err = core.WriteGeneratedCases(tmp, cases, true, w, pkgName, absRoot)
	if err != nil {
		return err
	}

	if hasMod {
		if err := core.TidyGoMod(tmp); err != nil {
			return err
		}
	}

	goBuildArgs := []string{"build", "-mod=mod"}
	if NeedsBuildVCSFlag(tmp) {
		goBuildArgs = append(goBuildArgs, "-buildvcs=false")
	}
	if opts.Verbose {
		goBuildArgs = append(goBuildArgs, "-v")
	}
	goBuildArgs = append(goBuildArgs, "./...")

	fmt.Fprintf(w, "cd %s && go %s\n\n", tmp, strings.Join(goBuildArgs, " "))

	goBuildCmd := exec.Command("go", goBuildArgs...)
	goBuildCmd.Dir = tmp

	if opts.Verbose {
		goBuildCmd.Stdout = w
		goBuildCmd.Stderr = w
		if err := goBuildCmd.Run(); err != nil {
			return fmt.Errorf("go build failed: %v", err)
		}
	} else {
		if out, err := goBuildCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("go build failed: %v\n%s", err, string(out))
		}
	}
	return nil
}
