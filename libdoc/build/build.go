package build

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/pathfmt"
)

func Build(dir string, opts core.Options) error {
	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}

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

	ctx, err := newGenerateContext(dir, opts, cases, w, true, opts.Verbose)
	if err != nil {
		return err
	}
	defer ctx.Close()

	ctx.announceRoots()

	if opts.Verbose {
		fmt.Fprintf(w, "doctest: %s\n\n", pathfmt.DisplayPath(dir))
		if _, err := core.DiscoverTreeCasesVerbose(dir, w); err != nil {
			return err
		}
		fmt.Fprintf(w, "─── %d test cases\n\n", len(cases))
	} else {
		fmt.Fprintf(w, "doctest: %s\n", pathfmt.DisplayPath(dir))
		fmt.Fprintf(w, "─── %d test cases\n", len(cases))
	}

	if err := ctx.writeCases(cases, true); err != nil {
		return err
	}

	goBuildArgs := []string{"build", "-mod=mod"}
	if NeedsBuildVCSFlag(ctx.genRoot) {
		goBuildArgs = append(goBuildArgs, "-buildvcs=false")
	}
	if opts.Verbose {
		goBuildArgs = append(goBuildArgs, "-v")
	}
	goBuildArgs = append(goBuildArgs, "./...")

	fmt.Fprintf(w, "cd %s && go %s\n\n", pathfmt.DisplayPath(ctx.genRoot), strings.Join(goBuildArgs, " "))

	goBuildCmd := exec.Command("go", goBuildArgs...)
	goBuildCmd.Dir = ctx.genRoot
	out, err := goBuildCmd.CombinedOutput()
	os.Stdout.Write(out)
	if err != nil {
		return fmt.Errorf("go build failed: %v", err)
	}

	if err := ctx.syncDump(); err != nil {
		return err
	}
	return nil
}