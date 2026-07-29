package build

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

type generateContext struct {
	w               io.Writer
	verbose         bool
	absRoot         string
	absModRoot      string
	modRoot         string
	modPath         string
	hasMod          bool
	genRoot         string
	compileRoot     string
	dumpDir         string
	internalCompile bool
	assertImport    bool
	assertCacheDir  string
	sessionImport   bool
	sessionCacheDir string
	modfilePath     string
	removeLegacyTmp bool
	// goCache is optional isolated GOCACHE for go mod tidy / go build children.
	goCache string
	// genBatch tracks -a wipe-once and desired emit paths (orphan prune).
	genBatch *core.GenBatch
	// generateOnly: defer orphan prune to workspace fan-in (multi-tree).
	generateOnly bool
	// forceA: wipe gen root once via genBatch before generate.
	forceA bool
	// unlockGenRoot releases LockGenRootWrites for shared mapping-gen safety.
	unlockGenRoot func()
	// unifiedMode: hierarchical ref packages + one suite package per DOCTEST tree.
	// Always true for normal generation; false only for internal-compile trees
	// (module-internal import path layout still uses classic AssembleTestSource).
	unifiedMode bool
	// subDir is the path-scope filter (opts.SubDir): abs or relative source path
	// under the DOCTEST tree (mid branch, leaf, or equal to tree root for full tree).
	subDir string
	closeOnce sync.Once
	// lifecycleMu serializes gen writes against interrupt cleanup. The SIGINT
	// handler acquires it and holds it through os.Exit so writeCases cannot
	// recreate .doctest_run_* / temp roots after RemoveAll.
	lifecycleMu sync.Mutex
	closed      bool
}

func newGenerateContext(dir string, opts core.Options, cases []core.TreeCase, w io.Writer, forBuild bool, verbose bool) (*generateContext, error) {
	absRoot, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	modRoot, modPath, hasMod := core.FindModuleRoot(absRoot)
	absModRoot, _ := core.MappingGenRoot(absRoot)
	// Always materialize assert-mod + session-mod replaces for external modules.
	// Multi-tree ./... shares one mapping-gen root per module; per-tree
	// assertImport=false would strip the assert replace and leave session
	// replace only — tidy then resolves assert via github.com/xhd2015/doctest@vX
	// (parent module also has session/) and fails with ambiguous session import.
	// Session is always required: assemble injects d *session.Doctest.
	assertImport := true
	sessionImport := true

	// Default generation is hierarchical unified (ref packages + suite).
	// Internal-compile trees keep classic AssembleTestSource (module-internal
	// import path layout).
	internalCompile := hasMod && core.CasesImportInternalPackage(cases, modPath)
	unifiedMode := !internalCompile

	ctx := &generateContext{
		w:               w,
		verbose:         verbose,
		absRoot:         absRoot,
		absModRoot:      absModRoot,
		modRoot:         modRoot,
		modPath:         modPath,
		hasMod:          hasMod,
		dumpDir:         opts.GenDir,
		internalCompile: internalCompile,
		assertImport:    assertImport,
		sessionImport:   sessionImport,
		goCache:         opts.GoCache,
		genBatch:        opts.GenBatch,
		generateOnly:    opts.GenerateOnly,
		forceA:          opts.ForceWithFlagA,
		unifiedMode:     unifiedMode,
		subDir:          opts.SubDir,
	}

	if assertImport {
		cacheDir, err := core.MaterializeAssertModule()
		if err != nil {
			return nil, err
		}
		ctx.assertCacheDir = cacheDir
	}
	if sessionImport {
		cacheDir, err := core.MaterializeSessionModule()
		if err != nil {
			return nil, err
		}
		ctx.sessionCacheDir = cacheDir
	}

	if ctx.internalCompile {
		compileRoot, err := core.NewInternalCompileRoot(modRoot)
		if err != nil {
			return nil, err
		}
		ctx.compileRoot = compileRoot
		ctx.genRoot = compileRoot
		return ctx, nil
	}

	genRoot := opts.GenDir
	if genRoot == "" {
		if forBuild {
			tmp, tmpErr := os.MkdirTemp("", "doctest-build-*")
			if tmpErr != nil {
				return nil, tmpErr
			}
			genRoot = tmp
			ctx.removeLegacyTmp = opts.RemoveTemp
		} else {
			var cacheErr error
			// Single warm mapping-gen path for hierarchical unified generation.
			genRoot, _, cacheErr = core.CacheMappingGenRoot(absRoot)
			if cacheErr != nil {
				return nil, cacheErr
			}
		}
	}
	if err := os.MkdirAll(genRoot, 0755); err != nil {
		return nil, err
	}
	// Serialize generate+prune per gen root so parallel nested suite leaves
	// sharing mapping-gen cannot clobber each other's packages / emit notes.
	ctx.unlockGenRoot = core.LockGenRootWrites(genRoot)
	// GenBatch: optional -a wipe-once; attach for emit-set orphan recording.
	if ctx.genBatch != nil {
		ctx.genBatch.Attach(genRoot)
		if ctx.forceA {
			if err := ctx.genBatch.WipeOnce(genRoot); err != nil {
				ctx.unlockGenRoot()
				ctx.unlockGenRoot = nil
				return nil, err
			}
		}
	} else if opts.ForceWithFlagA {
		// Library -a without shared GenBatch: still wipe this gen root.
		if err := core.WipeGenRoot(genRoot); err != nil {
			ctx.unlockGenRoot()
			ctx.unlockGenRoot = nil
			return nil, err
		}
	}
	ctx.genRoot = genRoot
	return ctx, nil
}

// removeTempsLocked deletes interrupt-scoped temps. Caller must hold lifecycleMu.
func (ctx *generateContext) removeTempsLocked() {
	if ctx.modfilePath != "" {
		// Internal-compile writes <modRoot>/.doctest.mod and runs go with
		// -modfile=…. Go's sum companion is the same path with .mod → .sum
		// (e.g. .doctest.sum). Remove both so consumer module roots stay clean.
		os.Remove(ctx.modfilePath)
		if sumPath := modfileCompanionSum(ctx.modfilePath); sumPath != "" {
			os.Remove(sumPath)
		}
	}
	if ctx.compileRoot != "" {
		os.RemoveAll(ctx.compileRoot)
	}
	if ctx.removeLegacyTmp && ctx.dumpDir == "" {
		os.RemoveAll(ctx.genRoot)
	}
}

// modfileCompanionSum returns the go.sum path that accompanies a -modfile path.
// Go derives it by replacing a trailing ".mod" with ".sum" (see cmd/go modload).
func modfileCompanionSum(modfilePath string) string {
	if !strings.HasSuffix(modfilePath, ".mod") {
		return ""
	}
	return strings.TrimSuffix(modfilePath, ".mod") + ".sum"
}

func (ctx *generateContext) Close() {
	ctx.closeOnce.Do(func() {
		ctx.lifecycleMu.Lock()
		defer ctx.lifecycleMu.Unlock()
		ctx.closed = true
		// Safety if writeCases returned early before finishGenOrphans.
		ctx.releaseGenWrite()
		ctx.removeTempsLocked()
	})
}

// withGenLock runs fn while holding lifecycleMu. Returns an error if generation
// was already closed (e.g. SIGINT cleanup started).
func (ctx *generateContext) withGenLock(fn func() error) error {
	ctx.lifecycleMu.Lock()
	defer ctx.lifecycleMu.Unlock()
	if ctx.closed {
		return fmt.Errorf("doctest: interrupted")
	}
	return fn()
}

func (ctx *generateContext) installInterruptCleanup() {
	if ctx.compileRoot == "" && !ctx.removeLegacyTmp {
		return
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		// Hold lifecycleMu until process exit so concurrent writeCases cannot
		// MkdirAll/WriteFile a removed compile temp back into existence.
		ctx.lifecycleMu.Lock()
		ctx.closeOnce.Do(func() {
			ctx.closed = true
			ctx.removeTempsLocked()
		})
		// Re-remove even if Close already ran without this lock held for exit.
		ctx.removeTempsLocked()
		os.Exit(130)
	}()
}

func (ctx *generateContext) announceRoots() {
	if ctx.dumpDir != "" {
		fmt.Fprintf(ctx.w, "→ %s\n\n", pathfmt.Short(ctx.dumpDir))
		return
	}
	fmt.Fprintf(ctx.w, "→ %s\n\n", pathfmt.Short(ctx.genRoot))
}

// writeCases generates packages for one tree. Shared gen-root bookkeeping
// (go.mod, tidy-done, doctest.gen-manifest) is serialized inside core via
// genModMu; tree-local package dirs may be written in parallel by multi-tree
// ./... prepare (no global writeCases lock).
//
// Ephemeral compile/build temps are written under lifecycleMu so SIGINT cleanup
// can RemoveAll without racing recreating writers.
func (ctx *generateContext) writeCases(cases []core.TreeCase, compileOnly bool) error {
	pkgName := "testcase"
	srcDir, origPkg, hasPkgUnderTest := core.ResolvePkgUnderTest(ctx.absRoot)
	if hasPkgUnderTest {
		pkgName = origPkg + "_tc"
	}

	if err := ctx.withGenLock(func() error {
		if !ctx.internalCompile {
			if err := core.WriteGoMod(ctx.genRoot, ctx.absModRoot, ctx.modPath, ctx.hasMod, ctx.assertImport, ctx.assertCacheDir, ctx.sessionImport, ctx.sessionCacheDir); err != nil {
				return err
			}
			if ctx.verbose && ctx.w != nil {
				fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(filepath.Join(ctx.genRoot, "go.mod")))
			}
		} else if ctx.assertImport || ctx.sessionImport {
			modfilePath, err := core.WriteInternalModfile(ctx.modRoot, ctx.assertCacheDir, ctx.sessionCacheDir)
			if err != nil {
				return err
			}
			ctx.modfilePath = modfilePath
			if ctx.verbose && ctx.w != nil {
				fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(modfilePath))
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if ctx.unifiedMode {
		return ctx.withGenLock(func() error {
			if err := ctx.writeUnifiedCases(cases, compileOnly, pkgName, hasPkgUnderTest, srcDir, origPkg); err != nil {
				return err
			}
			if err := core.FlushGenManifest(ctx.genRoot); err != nil {
				return err
			}
			return ctx.finishGenOrphans()
		})
	}

	// Internal-compile only: classic full-inline AssembleTestSource per leaf.
	// Per-leaf lock so SIGINT can clean up between leaves (and after a leaf finishes).
	for _, tc := range cases {
		tc := tc
		if err := ctx.withGenLock(func() error {
			absLeafDir := filepath.Join(ctx.absRoot, tc.Path)
			leafDir, err := core.GenDirForLeaf(ctx.genRoot, ctx.absModRoot, absLeafDir)
			if err != nil {
				return fmt.Errorf("gen dir for leaf %s: %w", tc.Path, err)
			}

			if hasPkgUnderTest {
				if _, err := core.CopySourceFiles(leafDir, srcDir, origPkg); err != nil {
					return fmt.Errorf("copy source files to %s: %w", leafDir, err)
				}
			}

			testPath, _, err := core.WriteGeneratedCase(leafDir, tc, compileOnly, pkgName, ctx.absRoot)
			if err != nil {
				return err
			}
			if ctx.verbose && ctx.w != nil {
				if compileOnly {
					fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(leafDir))
				} else {
					fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(testPath))
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return ctx.withGenLock(func() error {
		if ctx.hasMod && !ctx.internalCompile {
			if err := core.CondTidyGoMod(ctx.genRoot, ctx.goCache); err != nil {
				return err
			}
		}
		if err := core.FlushGenManifest(ctx.genRoot); err != nil {
			return err
		}
		return ctx.finishGenOrphans()
	})
}

// finishGenOrphans reconciles throwaway gen for this tree only (treeRel scope).
// Deferred when GenerateOnly (multi-tree workspace prunes later).
// Never full-gen-prunes: shared mapping-gen holds many trees / nested suite leaves.
// Releases gen-root write lock after prune so go test does not hold it.
//
// Path-scoped runs (SubDir under the tree) prune only under suiteRel (e.g.
// tree/mid) so sibling packages outside the user path are never orphan-deleted.
func (ctx *generateContext) finishGenOrphans() error {
	defer ctx.releaseGenWrite()
	if ctx.internalCompile || ctx.generateOnly || ctx.genBatch == nil {
		return nil
	}
	// Single-tree path: no sibling treeRels to exclude.
	// Path-scoped: prune only under the selected prefix (suiteRel).
	return ctx.genBatch.PruneTree(ctx.genRoot, ctx.suiteRel(), nil)
}

// suiteRel is the gen-relative placement for suite/registry/allleaves.
// Full tree (no SubDir, or SubDir == tree root) → treeRel.
// Mid/leaf path scope → path under module (e.g. tree/mid).
func (ctx *generateContext) suiteRel() string {
	treeRel := ctx.treeRel()
	if ctx.subDir == "" {
		return treeRel
	}
	subAbs := ctx.subDir
	if !filepath.IsAbs(subAbs) {
		subAbs = filepath.Join(ctx.absRoot, subAbs)
	}
	subAbs = filepath.Clean(subAbs)
	if subAbs == filepath.Clean(ctx.absRoot) {
		return treeRel
	}
	rel, err := filepath.Rel(ctx.absModRoot, subAbs)
	if err != nil || rel == "" || rel == "." {
		return treeRel
	}
	// Stay inside the tree: never place suite outside treeRel.
	treeClean := filepath.ToSlash(filepath.Clean(treeRel))
	relSlash := filepath.ToSlash(filepath.Clean(rel))
	if treeClean != "." && relSlash != treeClean && !strings.HasPrefix(relSlash, treeClean+"/") {
		return treeRel
	}
	return rel
}

// isPathScoped is true when the user path is a proper sub-prefix of the DOCTEST tree
// (mid branch or leaf), not the whole tree.
func (ctx *generateContext) isPathScoped() bool {
	return filepath.ToSlash(filepath.Clean(ctx.suiteRel())) != filepath.ToSlash(filepath.Clean(ctx.treeRel()))
}

// releaseGenWrite detaches the batch and unlocks gen-root serialize lock.
func (ctx *generateContext) releaseGenWrite() {
	if ctx.genBatch != nil && ctx.genRoot != "" {
		ctx.genBatch.Detach(ctx.genRoot)
	}
	if ctx.unlockGenRoot != nil {
		ctx.unlockGenRoot()
		ctx.unlockGenRoot = nil
	}
}

// treeRel is the doctest root path relative to the module root (or ".").
func (ctx *generateContext) treeRel() string {
	treeRel := "."
	if rel, relErr := filepath.Rel(ctx.absModRoot, ctx.absRoot); relErr == nil {
		treeRel = rel
	}
	return treeRel
}

// writeUnifiedCases generates ref root + intermediate packages once + registry +
// non-test leaf packages + __allleaves blank-import fan-in + suite iterator
// (one go test package/binary).
func (ctx *generateContext) writeUnifiedCases(cases []core.TreeCase, compileOnly bool, pkgName string, hasPkgUnderTest bool, srcDir, origPkg string) error {
	if len(cases) == 0 {
		return nil
	}

	rootDocs, _ := core.SplitRefSetupDocs(cases[0].SetupFiles)
	if len(rootDocs) == 0 {
		rootDocs = cases[0].SetupFiles
	}

	treeRel := ctx.treeRel()
	// Path-scoped mid/leaf: suite + registry + allleaves live under suiteRel
	// (e.g. tree/mid) so we do not rewrite tree-wide suite outside the path.
	// droot stays tree-scoped (shared types/Run for the DOCTEST root).
	suiteRel := ctx.suiteRel()
	rootDir := core.RefRootDirForTree(ctx.genRoot, treeRel)
	rootImport := core.RefRootImportForTree(treeRel)
	registryImport := core.UnifiedRegistryImportForTree(suiteRel)
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return err
	}

	// Package-under-test sources must live with Run/types so unexported
	// symbols remain callable (classic single-package semantics).
	rootPkgName := core.RefRootPkgName
	if hasPkgUnderTest {
		rootPkgName = pkgName
		if _, err := core.CopySourceFiles(rootDir, srcDir, origPkg); err != nil {
			return fmt.Errorf("copy source files to ref root: %w", err)
		}
	}
	rootSrc, err := core.AssembleRefRootSource(rootDocs, rootPkgName)
	if err != nil {
		return err
	}
	rootPath := filepath.Join(rootDir, "droot.go")
	if err := core.WriteFormattedGo(rootPath, rootSrc); err != nil {
		return fmt.Errorf("write unified ref root package: %w", err)
	}
	if ctx.verbose && ctx.w != nil {
		fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(rootPath))
	}

	// Intermediate packages once (parents first), shared across leaves — same as ref.
	if _, err := core.WriteRefIntermediatePackages(ctx.genRoot, treeRel, rootImport, rootDocs, cases); err != nil {
		return err
	}
	if ctx.verbose && ctx.w != nil {
		for _, g := range core.CollectUniqueRefIntermediates(cases) {
			interPath := filepath.Join(core.RefIntermediateDirForTree(ctx.genRoot, treeRel, g.Dir), core.RefIntermediateFileName)
			fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(interPath))
		}
	}

	leafImports := make([]string, 0, len(cases))
	for _, tc := range cases {
		absLeafDir := filepath.Join(ctx.absRoot, tc.Path)
		leafDir, err := core.GenDirForLeaf(ctx.genRoot, ctx.absModRoot, absLeafDir)
		if err != nil {
			return fmt.Errorf("gen dir for leaf %s: %w", tc.Path, err)
		}

		if hasPkgUnderTest {
			if _, err := core.CopySourceFiles(leafDir, srcDir, origPkg); err != nil {
				return fmt.Errorf("copy source files to %s: %w", leafDir, err)
			}
		}

		leafPath, _, err := core.WriteUnifiedLeafCase(leafDir, tc, compileOnly, pkgName, ctx.absRoot, rootImport, registryImport)
		if err != nil {
			return err
		}
		if ctx.verbose && ctx.w != nil {
			fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(leafPath))
		}

		leafRel, relErr := filepath.Rel(ctx.genRoot, leafDir)
		if relErr != nil {
			return fmt.Errorf("leaf import path for %s: %w", tc.Path, relErr)
		}
		leafImports = append(leafImports, core.LeafImportForTree(leafRel))
	}

	// Suite/registry/allleaves under suiteRel (path-local when path-scoped).
	if err := core.WriteUnifiedTreeExtras(ctx.genRoot, suiteRel, ctx.absRoot, leafImports); err != nil {
		return err
	}
	// Workspace __wreg is tree-wide fan-in; skip on path-scoped runs so we do
	// not rewrite packages outside the user path. Full-tree runs still write it.
	if !ctx.isPathScoped() {
		if err := core.WriteTreeWreg(ctx.genRoot, treeRel, ctx.absRoot); err != nil {
			return fmt.Errorf("write tree wreg: %w", err)
		}
	}
	if ctx.verbose && ctx.w != nil {
		fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(core.UnifiedRegistryDirForTree(ctx.genRoot, suiteRel)))
		fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(core.UnifiedAllLeavesDirForTree(ctx.genRoot, suiteRel)))
		fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(core.UnifiedSuiteDirForTree(ctx.genRoot, suiteRel)))
		if !ctx.isPathScoped() {
			fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(core.TreeWregDirForTree(ctx.genRoot, treeRel)))
		}
	}

	if ctx.hasMod && !ctx.internalCompile {
		if err := core.CondTidyGoMod(ctx.genRoot, ctx.goCache); err != nil {
			return err
		}
	}
	return core.FlushGenManifest(ctx.genRoot)
}

func (ctx *generateContext) syncDump() error {
	if !ctx.internalCompile || ctx.dumpDir == "" {
		return nil
	}
	if err := os.MkdirAll(ctx.dumpDir, 0755); err != nil {
		return err
	}
	return core.CopyGeneratedTree(ctx.compileRoot, ctx.dumpDir)
}

func (ctx *generateContext) runDir(absRoot string, opts core.Options, cases []core.TreeCase) (string, bool) {
	runDir := ctx.genRoot
	isSingleLeaf := false
	if opts.SubDir != "" {
		subDirAbs := opts.SubDir
		if !filepath.IsAbs(subDirAbs) {
			subDirAbs = filepath.Join(absRoot, subDirAbs)
		}
		if _, err := os.Stat(filepath.Join(subDirAbs, "ASSERT.md")); err == nil {
			isSingleLeaf = true
		}
		relSubDir, err := filepath.Rel(ctx.absModRoot, subDirAbs)
		if err == nil && relSubDir != "." {
			runDir = filepath.Join(ctx.genRoot, relSubDir)
		}
	}
	if !isSingleLeaf && len(cases) == 1 && cases[0].Path != "" {
		leafDir, _ := core.GenDirForLeaf(ctx.genRoot, ctx.absModRoot, filepath.Join(absRoot, cases[0].Path))
		runDir = leafDir
		isSingleLeaf = true
	}
	return runDir, isSingleLeaf
}

func (ctx *generateContext) scopedMultiRunDir(absRoot string) string {
	runDir := ctx.genRoot
	relRoot, err := filepath.Rel(ctx.absModRoot, absRoot)
	if err == nil && relRoot != "." {
		runDir = filepath.Join(ctx.genRoot, relRoot)
	}
	return runDir
}

func (ctx *generateContext) packageArgsForCases(runDir, absRoot string, cases []core.TreeCase) ([]string, error) {
	if ctx.unifiedMode {
		// Single suite package under suiteRel (tree-wide or path-local mid/leaf).
		suiteDir := core.UnifiedSuiteDirForTree(ctx.genRoot, ctx.suiteRel())
		rel, err := filepath.Rel(runDir, suiteDir)
		if err != nil {
			return nil, fmt.Errorf("package path for suite: %w", err)
		}
		if rel == "." {
			return []string{"."}, nil
		}
		return []string{"./" + filepath.ToSlash(rel)}, nil
	}
	seen := make(map[string]bool)
	args := make([]string, 0, len(cases))
	for _, tc := range cases {
		absLeaf := filepath.Join(absRoot, tc.Path)
		leafGen, err := core.GenDirForLeaf(ctx.genRoot, ctx.absModRoot, absLeaf)
		if err != nil {
			return nil, err
		}
		rel, err := filepath.Rel(runDir, leafGen)
		if err != nil {
			return nil, fmt.Errorf("package path for %s: %w", tc.Path, err)
		}
		arg := "./" + filepath.ToSlash(rel)
		if seen[arg] {
			continue
		}
		seen[arg] = true
		args = append(args, arg)
	}
	sort.Strings(args)
	return args, nil
}

func (ctx *generateContext) goCommandExtraArgs() []string {
	if ctx.internalCompile && (ctx.assertImport || ctx.sessionImport) && ctx.modfilePath != "" {
		return []string{"-modfile=" + ctx.modfilePath}
	}
	return nil
}
