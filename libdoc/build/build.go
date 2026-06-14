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

	absRoot, _ := filepath.Abs(dir)

	mappingGenRoot := opts.GenDir
	removeTmp := false
	if mappingGenRoot == "" {
		mappingGenRoot, err = os.MkdirTemp("", "doctest-build-*")
		if err != nil {
			return err
		}
		removeTmp = opts.RemoveTemp
	} else if err := os.MkdirAll(mappingGenRoot, 0755); err != nil {
		return err
	}

	if opts.Stderr != nil {
		fmt.Fprintf(opts.Stderr, "→ %s\n\n", mappingGenRoot)
	} else {
		fmt.Fprintf(os.Stderr, "→ %s\n\n", mappingGenRoot)
	}

	if removeTmp {
		defer os.RemoveAll(mappingGenRoot)
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

		_, err = core.WriteGeneratedCase(leafDir, tc, true, pkgName, absRoot)
		if err != nil {
			return err
		}
		if opts.Verbose {
			fmt.Fprintf(w, "→ %s\n", leafDir)
		}
	}

	if hasMod {
		if err := core.CondTidyGoMod(mappingGenRoot); err != nil {
			return err
		}
	}

	goBuildArgs := []string{"build", "-mod=mod"}
	if NeedsBuildVCSFlag(mappingGenRoot) {
		goBuildArgs = append(goBuildArgs, "-buildvcs=false")
	}
	if opts.Verbose {
		goBuildArgs = append(goBuildArgs, "-v")
	}
	goBuildArgs = append(goBuildArgs, "./...")

	fmt.Fprintf(w, "cd %s && go %s\n\n", mappingGenRoot, strings.Join(goBuildArgs, " "))

	goBuildCmd := exec.Command("go", goBuildArgs...)
	goBuildCmd.Dir = mappingGenRoot
	out, err := goBuildCmd.CombinedOutput()
	os.Stdout.Write(out)
	if err != nil {
		return fmt.Errorf("go build failed: %v", err)
	}
	return nil
}
