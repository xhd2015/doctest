package build

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xhd2015/doctest/libdoc/core"
)

func sha256HexOf(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func Test(dir string, opts core.Options) error {
	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}

	var cases []core.TreeCase
	var err error

	absRoot, _ := filepath.Abs(dir)

	tmp := opts.GenDir
	if tmp == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return err
		}
		h := sha256HexOf(absRoot)
		tmp = filepath.Join(cacheDir, "doctest", h)
	}
	if err := os.MkdirAll(tmp, 0755); err != nil {
		return err
	}

	fmt.Fprintf(w, "→ %s\n\n", tmp)

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

	sourceHash, err := core.ComputeSourceHash(dir, cases)
	if err != nil {
		return fmt.Errorf("compute source hash: %w", err)
	}

	hashMatch := false
	if existingHash, err := core.ReadHashFile(tmp); err == nil && existingHash == sourceHash {
		hashMatch = true
	}

	if !hashMatch {
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

		_, err = core.WriteGeneratedCases(tmp, cases, false, nil, pkgName, absRoot)
		if err != nil {
			return err
		}

		if hasMod {
			if err := core.TidyGoMod(tmp); err != nil {
				return err
			}
		}

		if err := core.WriteHashFile(tmp, sourceHash); err != nil {
			return err
		}
	}

	testArgs := []string{"test", "-mod=mod", "-v"}
	if NeedsBuildVCSFlag(tmp) {
		testArgs = append(testArgs, "-buildvcs=false")
	}
	if opts.Count > 0 {
		testArgs = append(testArgs, fmt.Sprintf("-count=%d", opts.Count))
	}
	if len(cases) > 0 {
		var runPattern strings.Builder
		runPattern.WriteString("^(")
		for i, tc := range cases {
			if i > 0 {
				runPattern.WriteByte('|')
			}
			runPattern.WriteString(core.TestFuncName(tc))
		}
		runPattern.WriteString(")$")
		testArgs = append(testArgs, "-run", runPattern.String())
	}
	testArgs = append(testArgs, ".")

	fmt.Fprintf(w, "cd %s && go %s\n\n", tmp, strings.Join(testArgs, " "))

	goTestCmd := exec.Command("go", testArgs...)
	goTestCmd.Dir = tmp
	goTestCmd.Stdout = w
	goTestCmd.Stderr = w
	if err := goTestCmd.Run(); err != nil {
		return fmt.Errorf("go test failed: %v", err)
	}
	return nil
}
