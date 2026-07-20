package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

// TreePrep is the result of generating one DOCTEST root for a workspace run.
type TreePrep struct {
	AbsRoot string
	TreeRel string
	GenRoot string
	Unified bool
	Cases   []core.TreeCase
	Skipped []core.SkippedCase
	Stats   TestRunStats // phases + totals from generate-only
}

// PrepareTree discovers, filters, and generates one tree without go test.
// Unified trees write __wreg for later workspace fan-in.
func PrepareTree(dir string, opts core.Options) (TreePrep, error) {
	o := opts
	o.GenerateOnly = true
	o.SuppressResultSummary = true
	stats, err := TestWithStats(dir, o)
	prep := TreePrep{
		AbsRoot: stats.AbsRoot,
		TreeRel: stats.TreeRel,
		GenRoot: stats.GenRoot,
		Unified: stats.Unified,
		Skipped: stats.Skipped,
		Stats:   stats,
	}
	if prep.AbsRoot == "" {
		if abs, aerr := filepath.Abs(dir); aerr == nil {
			prep.AbsRoot = abs
		} else {
			prep.AbsRoot = dir
		}
	}
	// Re-discover cases for metrics/leaf attribution (same filters as TestWithStats).
	if cases, derr := core.DiscoverTreeCases(dir); derr == nil {
		if o.SubDir != "" {
			cases = core.FilterBySubDir(cases, dir, o.SubDir)
		}
		if o.ChangedOnly {
			gitRoot, changedFiles, cerr := core.ChangedGitFiles(dir)
			if cerr == nil {
				cases = core.FilterByChangedFiles(cases, dir, gitRoot, changedFiles)
			}
		}
		cases, _ = core.FilterCasesByLabel(cases, o)
		prep.Cases = cases
	}
	if err != nil {
		return prep, err
	}
	return prep, nil
}

// RunWorkspace writes __workspace fan-in for the given unified preps and runs
// a single go test on __workspace/suite.
func RunWorkspace(preps []TreePrep, opts core.Options) (TestRunStats, error) {
	var stats TestRunStats
	if len(preps) == 0 {
		return stats, fmt.Errorf("workspace: no trees to run")
	}
	genRoot := preps[0].GenRoot
	if genRoot == "" {
		return stats, fmt.Errorf("workspace: empty gen root")
	}
	for _, p := range preps {
		if !p.Unified {
			return stats, fmt.Errorf("workspace: tree %s is not unified", p.AbsRoot)
		}
		if p.GenRoot != "" && filepath.Clean(p.GenRoot) != filepath.Clean(genRoot) {
			return stats, fmt.Errorf("workspace: mixed gen roots %s vs %s", genRoot, p.GenRoot)
		}
		stats.Total += p.Stats.Total
		stats.Skipped = append(stats.Skipped, p.Skipped...)
		stats.Phases = append(stats.Phases, p.Stats.Phases...)
	}
	if stats.Total == 0 {
		return stats, nil
	}

	treeRels := make([]string, 0, len(preps))
	seen := map[string]bool{}
	for _, p := range preps {
		if p.Stats.Total == 0 {
			continue
		}
		tr := p.TreeRel
		if tr == "" {
			tr = "."
		}
		if seen[tr] {
			continue
		}
		seen[tr] = true
		treeRels = append(treeRels, tr)
	}
	sort.Slice(treeRels, func(i, j int) bool {
		return filepath.ToSlash(treeRels[i]) < filepath.ToSlash(treeRels[j])
	})

	if err := core.WriteWorkspaceExtras(genRoot, treeRels); err != nil {
		return stats, err
	}
	if err := core.CondTidyGoMod(genRoot); err != nil {
		return stats, err
	}

	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	suiteDir := core.WorkspaceSuiteDir(genRoot)
	rel, err := filepath.Rel(genRoot, suiteDir)
	if err != nil {
		return stats, err
	}
	packageArgs := []string{"."}
	if rel != "." {
		packageArgs = []string{"./" + filepath.ToSlash(rel)}
	}

	flagArgs := []string{"test", "-mod=mod"}
	if opts.Verbose {
		flagArgs = append(flagArgs, "-v")
	}
	if NeedsBuildVCSFlag(genRoot) {
		flagArgs = append(flagArgs, "-buildvcs=false")
	}
	if opts.Count > 0 {
		flagArgs = append(flagArgs, fmt.Sprintf("-count=%d", opts.Count))
	}
	if opts.Timeout > 0 {
		flagArgs = append(flagArgs, fmt.Sprintf("-timeout=%s", opts.Timeout))
	}
	if opts.CPUProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-cpuprofile=%s", opts.CPUProfile))
	}
	if opts.MemProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-memprofile=%s", opts.MemProfile))
	}
	if opts.MemProfileRate != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-memprofilerate=%d", *opts.MemProfileRate))
	}
	if opts.BlockProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-blockprofile=%s", opts.BlockProfile))
	}
	if opts.BlockProfileRate != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-blockprofilerate=%d", *opts.BlockProfileRate))
	}
	if opts.MutexProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-mutexprofile=%s", opts.MutexProfile))
	}
	if opts.MutexProfileFraction != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-mutexprofilefraction=%d", *opts.MutexProfileFraction))
	}
	if opts.Trace != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-trace=%s", opts.Trace))
	}
	if opts.OutputDir != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-outputdir=%s", opts.OutputDir))
	}
	if opts.CoverProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-coverprofile=%s", opts.CoverProfile))
	}
	if opts.Cover {
		flagArgs = append(flagArgs, "-cover")
	}

	displayArgs := displayGoArgs(append(append([]string(nil), flagArgs...), packageArgs...))
	if opts.Verbose {
		fmt.Fprintf(w, "cd %s && go %s\n\n", pathfmt.Short(genRoot), strings.Join(displayArgs, " "))
	} else {
		fmt.Fprintf(w, "doctest: workspace (%d trees, %d tests)\n", len(treeRels), stats.Total)
		fmt.Fprintf(w, "cd %s && go %s\n", pathfmt.Short(genRoot), strings.Join(displayArgs, " "))
	}

	sessionID := core.DoctestSessionIDForRun()
	goCache := opts.GoCache
	style := newColorStyle(opts.Color, stdout)

	tGo := time.Now()
	var runErr error
	if opts.Verbose {
		execArgs := append(append([]string(nil), flagArgs...), packageArgs...)
		goTestCmd := exec.Command("go", execArgs...)
		goTestCmd.Dir = genRoot
		goTestCmd.Env = goTestEnvFull(sessionID, goCache, opts.MetricsNestSink)
		out, err := goTestCmd.CombinedOutput()
		stdout.Write(out)
		runErr = err
		stats.Passed = passedCases(stats.Total, countFailuresFromGoTestOutput(out))
	} else {
		result, err := runGoTestJSONOnce(genRoot, append(append([]string(nil), flagArgs...), packageArgs...), sessionID, goCache, opts.MetricsNestSink, stdout, style)
		runErr = err
		goTestElapsed := time.Since(tGo)
		stats.Phases = append(stats.Phases, PhaseTiming{Name: "go_test", ElapsedNs: goTestElapsed.Nanoseconds()})
		if result.passCount+result.failCount > 0 {
			stats.Passed = result.passCount
			if result.passCount+result.failCount != stats.Total {
				stats.Total = result.passCount + result.failCount
			}
		} else {
			stats.Passed = passedCases(stats.Total, result.failCount)
		}
		fmt.Fprintln(stdout, formatSummary(style, result.passCount+result.failCount, result.passCount, result.failCount, result.cachedCount, goTestElapsed))
		for _, line := range result.failLines {
			fmt.Fprintln(stdout, line)
		}
		for _, line := range result.detailLines {
			fmt.Fprintln(stdout, line)
		}
		if len(result.stderrData) > 0 {
			stdout.Write(result.stderrData)
		}
		var allCases []core.TreeCase
		for _, p := range preps {
			allCases = append(allCases, p.Cases...)
		}
		stats.LeafTimings = leafTimingsFromSubtests(allCases, result, goTestElapsed)
	}

	if !opts.Verbose {
		fmt.Fprintln(w)
	}
	stats.Elapsed = time.Since(tGo)
	stats.GenRoot = genRoot
	stats.Unified = true

	if runErr != nil {
		return stats, fmt.Errorf("go test: %w", runErr)
	}
	return stats, nil
}
