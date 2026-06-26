package build

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

func Build(dir string, opts core.Options) error {
	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}

	var allCases, cases []core.TreeCase
	var err error

	allCases, err = core.DiscoverTreeCases(dir)
	if err != nil {
		return err
	}
	if opts.SubDir != "" {
		allCases = core.FilterBySubDir(allCases, dir, opts.SubDir)
	}

	cases = allCases
	if opts.ChangedOnly {
		gitRoot, changedFiles, err := core.ChangedGitFiles(dir)
		if err != nil {
			return err
		}
		changedInfo := core.ChangedRunInfoForTree(allCases, dir, gitRoot, changedFiles)
		cases = core.FilterByChangedFiles(allCases, dir, gitRoot, changedFiles)
		if core.ShouldAnnounceChangedRun(changedInfo, opts.Verbose) {
			fmt.Fprintln(w, core.FormatDoctestAnnouncement(pathfmt.Short(dir), changedInfo, true, 0))
		}
		if len(cases) == 0 {
			return nil
		}
	}
	if len(cases) == 0 {
		return fmt.Errorf("%s: no runnable test cases found", dir)
	}

	ctx, err := newGenerateContext(dir, opts, cases, w, true, opts.Verbose)
	if err != nil {
		return err
	}
	ctx.installInterruptCleanup()
	defer ctx.Close()

	ctx.announceRoots()

	if opts.Verbose {
		fmt.Fprintf(w, "doctest: %s\n\n", pathfmt.Short(dir))
		if _, err := core.DiscoverTreeCasesVerbose(dir, w); err != nil {
			return err
		}
		fmt.Fprintf(w, "─── %d test cases\n\n", len(cases))
	} else {
		fmt.Fprintf(w, "doctest: %s\n", pathfmt.Short(dir))
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

	fmt.Fprintf(w, "cd %s && go %s\n\n", pathfmt.Short(ctx.genRoot), strings.Join(displayGoArgs(goBuildArgs), " "))

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