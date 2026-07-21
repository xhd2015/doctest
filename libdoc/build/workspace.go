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
		Cases:   stats.Cases, // reuse filtered cases — avoid second full discover
		Stats:   stats,
	}
	if prep.AbsRoot == "" {
		if abs, aerr := filepath.Abs(dir); aerr == nil {
			prep.AbsRoot = abs
		} else {
			prep.AbsRoot = dir
		}
	}
	if err != nil {
		return prep, err
	}
	return prep, nil
}

// RunWorkspace writes __workspace fan-in for the given unified preps and runs
// a single go test. Same gen root → classic __workspace suite. Mixed gen roots
// (multi-module ./...) → toplevel/__hub go.mod + suite calling each RunAll.
func RunWorkspace(preps []TreePrep, opts core.Options) (TestRunStats, error) {
	var stats TestRunStats
	if len(preps) == 0 {
		return stats, fmt.Errorf("workspace: no trees to run")
	}

	active := make([]TreePrep, 0, len(preps))
	for _, p := range preps {
		if !p.Unified {
			return stats, fmt.Errorf("workspace: tree %s is not unified", p.AbsRoot)
		}
		if p.GenRoot == "" {
			return stats, fmt.Errorf("workspace: empty gen root for %s", p.AbsRoot)
		}
		stats.Skipped = append(stats.Skipped, p.Skipped...)
		stats.Phases = append(stats.Phases, p.Stats.Phases...)
		if p.Stats.Total == 0 {
			continue
		}
		active = append(active, p)
		stats.Total += p.Stats.Total
	}
	if stats.Total == 0 {
		return stats, nil
	}

	// Group by gen root.
	byGen := map[string][]TreePrep{}
	var genOrder []string
	for _, p := range active {
		g := filepath.Clean(p.GenRoot)
		if _, ok := byGen[g]; !ok {
			genOrder = append(genOrder, g)
		}
		byGen[g] = append(byGen[g], p)
	}
	sort.Strings(genOrder)

	if len(genOrder) == 1 {
		return runWorkspaceSingleGen(active, genOrder[0], stats, opts)
	}
	return runWorkspaceMultiModHub(active, byGen, genOrder, stats, opts)
}

func runWorkspaceSingleGen(preps []TreePrep, genRoot string, stats TestRunStats, opts core.Options) (TestRunStats, error) {
	treeRels := make([]string, 0, len(preps))
	seen := map[string]bool{}
	for _, p := range preps {
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

	suiteDir := core.WorkspaceSuiteDir(genRoot)
	rel, err := filepath.Rel(genRoot, suiteDir)
	if err != nil {
		return stats, err
	}
	packageArgs := []string{"."}
	if rel != "." {
		packageArgs = []string{"./" + filepath.ToSlash(rel)}
	}

	return finishWorkspaceGoTest(preps, genRoot, genRoot, packageArgs, len(treeRels), stats, opts)
}

func runWorkspaceMultiModHub(preps []TreePrep, byGen map[string][]TreePrep, genOrder []string, stats TestRunStats, opts core.Options) (TestRunStats, error) {
	// Per gen root: optional multi-tree workspace extras, unique module path, suite import.
	replaceByMod := map[string]string{} // module path → gen root abs
	var members []memberSuite
	usedAlias := map[string]bool{}

	for _, genRoot := range genOrder {
		group := byGen[genRoot]
		treeRels := make([]string, 0, len(group))
		seen := map[string]bool{}
		for _, p := range group {
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
		multiTree := len(treeRels) > 1
		if multiTree {
			if err := core.WriteWorkspaceExtras(genRoot, treeRels); err != nil {
				return stats, err
			}
		}
		if err := core.CondTidyGoMod(genRoot); err != nil {
			return stats, err
		}
		unique := uniqueWorkModulePath(genRoot)
		if err := ensureWorkModulePath(genRoot, unique); err != nil {
			return stats, fmt.Errorf("workspace multi-mod: module path for %s: %w", genRoot, err)
		}
		replaceByMod[unique] = genRoot

		// One RunAll entry per gen root (workspace suite if multi-tree).
		var suiteImp, subName string
		if multiTree {
			suiteImp = suiteImportForPrep(unique, ".", true)
			subName = filepath.Base(genRoot)
		} else {
			tr := treeRels[0]
			suiteImp = suiteImportForPrep(unique, tr, false)
			subName = tr
			if subName == "." {
				subName = filepath.Base(genRoot)
			}
		}
		subName = filepath.ToSlash(subName)
		subName = strings.ReplaceAll(subName, "/", "__")
		members = append(members, memberSuite{
			Alias: aliasForImportPath(suiteImp, usedAlias),
			Path:  suiteImp,
			Name:  subName,
		})
	}

	toplevel := pickToplevelGenRoot(genOrder)
	if toplevel == "" {
		return stats, fmt.Errorf("workspace multi-mod: empty toplevel gen root")
	}
	// If toplevel is only a common parent without go.mod, still OK — hub lives under it.
	hubDir, err := writeMultiModHub(toplevel, members, replaceByMod)
	if err != nil {
		return stats, fmt.Errorf("workspace multi-mod hub: %w", err)
	}

	return finishWorkspaceGoTest(preps, hubDir, hubDir, []string{"./suite"}, len(preps), stats, opts)
}

// finishWorkspaceGoTest runs go test for workspace (single-gen or multi-mod hub).
// runDir is the process cwd (gen root or __hub).
// When opts.BypassGoTest is set, workspace/hub files are already written by the
// caller; this returns without exec'ing go test.
func finishWorkspaceGoTest(preps []TreePrep, runDir, genRootLabel string, packageArgs []string, treeCount int, stats TestRunStats, opts core.Options) (TestRunStats, error) {
	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	if opts.BypassGoTest {
		label := "workspace"
		if strings.Contains(runDir, HubDirName) {
			label = "workspace hub"
		}
		fmt.Fprintf(w, "doctest: %s (%d trees, %d tests)\n", label, treeCount, stats.Total)
		fmt.Fprintf(w, "doctest: DOCTEST_DEBUG bypass-go-test=1 (go test skipped after workspace write)\n")
		stats.GoTestBypassed = true
		stats.Passed = 0
		stats.GenRoot = genRootLabel
		stats.Unified = true
		stats.Phases = append(stats.Phases, PhaseTiming{Name: "go_test", ElapsedNs: 0})
		return stats, nil
	}

	flagArgs := []string{"test", "-mod=mod"}
	if opts.Verbose {
		flagArgs = append(flagArgs, "-v")
	}
	// buildvcs: check first gen root from preps
	checkRoot := preps[0].GenRoot
	if NeedsBuildVCSFlag(checkRoot) {
		flagArgs = append(flagArgs, "-buildvcs=false")
	}
	if opts.Count > 0 {
		flagArgs = append(flagArgs, fmt.Sprintf("-count=%d", opts.Count))
	}
	if opts.ForceWithFlagA {
		flagArgs = append(flagArgs, "-a")
	}
	// nil = omit (go default 10m); non-nil including 0 = pass -timeout=…
	if opts.Timeout != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-timeout=%s", *opts.Timeout))
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
		fmt.Fprintf(w, "cd %s && go %s\n\n", pathfmt.Short(runDir), strings.Join(displayArgs, " "))
	} else {
		label := "workspace"
		if strings.Contains(runDir, HubDirName) {
			label = "workspace hub"
		}
		fmt.Fprintf(w, "doctest: %s (%d trees, %d tests)\n", label, treeCount, stats.Total)
		fmt.Fprintf(w, "cd %s && go %s\n", pathfmt.Short(runDir), strings.Join(displayArgs, " "))
	}

	sessionID := core.DoctestSessionIDForRun()
	goCache := opts.GoCache
	style := newColorStyle(opts.Color, stdout)

	tGo := time.Now()
	var runErr error
	if opts.Verbose {
		execArgs := append(append([]string(nil), flagArgs...), packageArgs...)
		goTestCmd := exec.Command("go", execArgs...)
		goTestCmd.Dir = runDir
		goTestCmd.Env = goTestEnvFull(sessionID, goCache, opts.MetricsNestSink, "", "")
		out, err := goTestCmd.CombinedOutput()
		stdout.Write(out)
		runErr = err
		stats.Passed = passedCases(stats.Total, countFailuresFromGoTestOutput(out))
		if err != nil {
			if msg := goTestTimeoutErrorLine(string(out)); msg != "" {
				fmt.Fprintln(w, msg)
			}
		}
	} else {
		result, err := runGoTestJSONOnce(runDir, append(append([]string(nil), flagArgs...), packageArgs...), sessionID, goCache, opts.MetricsNestSink, "", "", stdout, style)
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
		printGoTestTimeoutError(w, stdout, result)
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
	stats.GenRoot = genRootLabel
	stats.Unified = true

	if runErr != nil {
		return stats, fmt.Errorf("go test: %w", runErr)
	}
	return stats, nil
}
