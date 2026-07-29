package build

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/leafcache"
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
	// Share GenBatch across multi-tree prepare when caller provided one; else
	// allocate so this prepare still records emit paths for a later prune.
	if o.GenBatch == nil {
		o.GenBatch = core.NewGenBatch()
	}
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
		// Planned accumulates discovery leaf counts across trees.
		planned := p.Stats.Planned
		if planned == 0 {
			planned = p.Stats.Total
		}
		stats.Planned += planned
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

	// Hold gen-root write lock only for extras/tidy/prune — not during go test
	// (nested suite leaves must not block on the outer suite's go test).
	if err := prepareWorkspaceGen(genRoot, treeRels, opts); err != nil {
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

// prepareWorkspaceGen writes hub packages, tidies, and prunes under genRoot.
// Caller must not hold gen-root lock across go test.
//
// Tree package prune is intentionally skipped here. Multi-tree GenerateOnly
// prepare + warm second runs (hash-hit writes) have left incomplete GenBatch
// emit notes; pruning then deleted live packages (e.g. tree-b/__wreg) while
// __alltrees still imported them. Single-tree finishGenOrphans still prunes.
// Full clean remains available via -a / --cold-cache. Hub-only prune keeps
// __workspace free of stale hub files.
func prepareWorkspaceGen(genRoot string, treeRels []string, opts core.Options) error {
	unlock := core.LockGenRootWrites(genRoot)
	defer unlock()
	if opts.GenBatch != nil {
		opts.GenBatch.Attach(genRoot)
		defer opts.GenBatch.Detach(genRoot)
	}
	if err := core.WriteWorkspaceExtras(genRoot, treeRels); err != nil {
		return err
	}
	if err := core.CondTidyGoMod(genRoot, opts.GoCache); err != nil {
		return err
	}
	if opts.GenBatch != nil {
		return opts.GenBatch.PruneWorkspace(genRoot)
	}
	return nil
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
			if err := prepareWorkspaceGen(genRoot, treeRels, opts); err != nil {
				return stats, err
			}
		} else {
			// Single tree under this gen root: tidy only (no tree prune here;
			// GenerateOnly already deferred; avoid warm-run emit gaps).
			unlock := core.LockGenRootWrites(genRoot)
			err := core.CondTidyGoMod(genRoot, opts.GoCache)
			unlock()
			if err != nil {
				return stats, err
			}
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
	hubDir, err := writeMultiModHub(toplevel, members, replaceByMod, opts.GoCache)
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
	if opts.Race {
		flagArgs = append(flagArgs, "-race")
	}

	goTestBin, err := resolveWorkspaceGoTestBinary(preps, opts)
	if err != nil {
		return stats, err
	}

	// Apply project test.config.json for xgo (first prep with a readable config;
	// multi-module workspaces typically share one consumer project).
	xgoApply, err := workspaceXgoTestConfigApply(goTestBin, preps)
	if err != nil {
		return stats, err
	}
	if len(xgoApply.Flags) > 0 {
		flagArgs = append(flagArgs, xgoApply.Flags...)
	}

	displayArgs := displayGoArgs(append(append([]string(nil), flagArgs...), packageArgs...))
	if len(xgoApply.ProgArgs) > 0 {
		displayArgs = append(displayArgs, "-args")
		displayArgs = append(displayArgs, xgoApply.ProgArgs...)
	}
	// Always print planned trees/tests before go test, including Verbose.
	label := "workspace"
	if strings.Contains(runDir, HubDirName) {
		label = "workspace hub"
	}
	fmt.Fprintf(w, "doctest: %s (%d trees, %d tests)\n", label, treeCount, stats.Total)
	if opts.Verbose {
		if xgoApply.ConfigPath != "" {
			fmt.Fprintf(w, "doctest: xgo config %s\n", pathfmt.Short(xgoApply.ConfigPath))
		}
		fmt.Fprintf(w, "cd %s && %s %s\n\n", pathfmt.Short(runDir), goTestBin, strings.Join(displayArgs, " "))
	} else {
		fmt.Fprintf(w, "cd %s && %s %s\n", pathfmt.Short(runDir), goTestBin, strings.Join(displayArgs, " "))
	}

	sessionID := core.SessionIDFromOpts(opts)
	goCache := opts.GoCache
	style := newColorStyle(opts.Color, stdout)

	// Multi-prep leaf-cache: PreparePassPlan over all trees, tree-qualified skip env.
	// leafKeys + skipPaths feed stream PutPass and grey warm-skip progress dots.
	leafKeys, skipPaths := prepareWorkspaceLeafCache(preps, opts)
	leafSkipEnv := leafcache.FormatSkipPaths(skipPaths)

	// Always use go test -json for Pass/Fail/Run suite accounting. Verbose is
	// presentation only (stream more Output events); same counts as quiet.
	tGo := time.Now()
	result, runErr := runGoTestJSONOnce(goTestBin, runDir, append(append([]string(nil), flagArgs...), packageArgs...), sessionID, goCache, opts.MetricsNestSink, "", leafSkipEnv, xgoApply.Env, xgoApply.ProgArgs, leafKeys, stdout, style, opts.Verbose)
	goTestElapsed := time.Since(tGo)
	stats.Phases = append(stats.Phases, PhaseTiming{Name: "go_test", ElapsedNs: goTestElapsed.Nanoseconds()})
	// Discovery planned count before Total is rewritten to actual_run.
	if stats.Planned == 0 {
		stats.Planned = stats.Total
	}
	if result.timeoutError != "" {
		stats.TimedOut = true
	}
	// Prefer JSON suite-leaf accounting. actual_run = pass+fail (exclude
	// runtime t.Skip from denominator). SkipCount is separate from label skips.
	// On timeout: never invent phantom passes from planned − failCount.
	// On build failed with no leaf events: honest 0 Run / 0 Cached.
	nCases := 0
	for _, p := range preps {
		nCases += len(p.Cases)
	}
	applyGoTestLeafStats(&stats, &result, nCases, skipPaths, opts)
	// Map suite fail bare leaf paths → FormatLeafIdentity for RecordPasses.
	failed := workspaceFailedIdentities(preps, result.suiteLeafFailed)
	recordLeafCachePasses(leafKeys, failed, runErr == nil && result.failCount == 0 && !result.buildFailed)

	// Quiet path: compact progress summary. Verbose already streamed Output events.
	// Print order: progress → build diagnostics → package FAIL → fail dumps →
	// Error/hint (PASS/FAIL is printed by the runner after return when SuppressResultSummary).
	if !opts.Verbose {
		fmt.Fprintln(stdout, formatSummary(style, result.passCount+result.failCount, result.passCount, result.failCount, result.cachedCount, goTestElapsed))
		for _, line := range result.buildOutputLines {
			fmt.Fprintln(stdout, line)
		}
		for _, line := range result.failLines {
			fmt.Fprintln(stdout, line)
		}
		for _, line := range result.detailLines {
			fmt.Fprintln(stdout, line)
		}
	} else if len(result.buildOutputLines) > 0 {
		for _, line := range result.buildOutputLines {
			fmt.Fprintln(stdout, line)
		}
	}
	if len(result.stderrData) > 0 {
		stdout.Write(result.stderrData)
	}
	printGoTestTimeoutError(stdout, result, style)
	captureVCSStatusFromBuffers(&result)
	printGoTestVCSStatusError(stdout, result, style)
	var allCases []core.TreeCase
	for _, p := range preps {
		allCases = append(allCases, p.Cases...)
	}
	stats.LeafTimings = leafTimingsFromSubtests(allCases, result, goTestElapsed)

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

// workspaceXgoTestConfigApply picks project test.config.json from the first
// prep that has one under its module root. Suite module is empty so include-as
// always applies when a project module path is known (workspace runs from gen/hub).
func workspaceXgoTestConfigApply(goTestBin string, preps []TreePrep) (core.XgoTestConfigApply, error) {
	if strings.TrimSpace(goTestBin) != "xgo" {
		return core.XgoTestConfigApply{}, nil
	}
	for _, p := range preps {
		modRoot, modPath, ok := core.FindModuleRoot(p.AbsRoot)
		if !ok || modRoot == "" {
			continue
		}
		if core.FindXgoTestConfigPath(modRoot) == "" {
			continue
		}
		// Empty suite path → always request include-as (gen/hub ≠ project).
		return core.BuildXgoTestConfigApply(goTestBin, modRoot, modPath, "")
	}
	// No config file: still try include-as from first prep module.
	for _, p := range preps {
		modRoot, modPath, ok := core.FindModuleRoot(p.AbsRoot)
		if !ok || modPath == "" {
			continue
		}
		return core.BuildXgoTestConfigApply(goTestBin, modRoot, modPath, "")
	}
	return core.XgoTestConfigApply{}, nil
}

// prepareWorkspaceLeafCache builds multi-tree PassPlan keys (FormatLeafIdentity)
// and env-safe skip tokens (FormatLeafIdentityEnv) for DOCTEST_LEAF_CACHE_SKIP_PATHS.
// Store I/O errors are ignored (best-effort; suite continues).
func prepareWorkspaceLeafCache(preps []TreePrep, opts core.Options) (keys map[string]string, skipPaths []string) {
	keys = make(map[string]string)
	var leaves []leafcache.LeafRef
	for _, p := range preps {
		root := p.AbsRoot
		for _, tc := range p.Cases {
			leaves = append(leaves, leafcache.LeafRef{TreeRoot: root, LeafRel: tc.Path})
		}
	}
	if len(leaves) == 0 {
		return keys, nil
	}
	storeRoot, err := leafcache.ResolveStoreRoot()
	if err != nil {
		return keys, nil
	}
	store, err := leafcache.NewStore(storeRoot)
	if err != nil {
		return keys, nil
	}
	goVer := runtime.Version()
	enabled := leafcache.SkipEnabled(opts.Count, opts.ForceWithFlagA)
	plan, err := leafcache.PreparePassPlan(store, leaves, goVer, enabled)
	if err != nil {
		return keys, nil
	}
	keys = plan.Keys
	for _, id := range plan.Skip {
		skipPaths = append(skipPaths, leafcache.IdentityToEnv(id))
	}
	// plan.Skip is already sorted; IdentityToEnv preserves order.
	return keys, skipPaths
}

// workspaceFailedIdentities maps suiteLeafFailed bare leaf paths to
// FormatLeafIdentity tokens for each matching prep case.
func workspaceFailedIdentities(preps []TreePrep, suiteLeafFailed map[string]bool) map[string]bool {
	failed := make(map[string]bool)
	if len(suiteLeafFailed) == 0 {
		return failed
	}
	for _, p := range preps {
		for _, tc := range p.Cases {
			if !suiteLeafFailed[tc.Path] {
				// Also accept slash-normalized path.
				if !suiteLeafFailed[filepath.ToSlash(tc.Path)] {
					continue
				}
			}
			failed[leafcache.FormatLeafIdentity(p.AbsRoot, tc.Path)] = true
		}
	}
	return failed
}
