package build

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
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
	dumpDir         string
	assertImport    bool
	assertCacheDir  string
	sessionImport   bool
	sessionCacheDir string
	// vendorBridges is current-generation metadata returned while writing the
	// generated go.mod. It is passed to pre_test overlay normalization rather
	// than rediscovered from a retained generated root.
	vendorBridges   []core.VendorBridgeMapping
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
	// subDir is the path-scope filter (opts.SubDir): abs or relative source path
	// under the DOCTEST tree (mid branch, leaf, or equal to tree root for full tree).
	subDir    string
	closeOnce sync.Once
	// lifecycleMu serializes gen writes against interrupt cleanup. The SIGINT
	// handler acquires it and holds it through os.Exit so writeCases cannot
	// recreate temp roots after RemoveAll.
	lifecycleMu sync.Mutex
	closed      bool
}

func newGenerateContext(dir string, opts core.Options, w io.Writer, forBuild bool, verbose bool) (*generateContext, error) {
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
	//
	// Always hierarchical unified (layout A + Kind A/B shims). Parent/product
	// internal imports use Kind B __doctest_internal_expose rewrite — never
	// classic multi-leaf under .doctest_run_*.
	assertImport := true
	sessionImport := true

	ctx := &generateContext{
		w:             w,
		verbose:       verbose,
		absRoot:       absRoot,
		absModRoot:    absModRoot,
		modRoot:       modRoot,
		modPath:       modPath,
		hasMod:        hasMod,
		dumpDir:       opts.GenDir,
		assertImport:  assertImport,
		sessionImport: sessionImport,
		goCache:       opts.GoCache,
		genBatch:      opts.GenBatch,
		generateOnly:  opts.GenerateOnly,
		forceA:        opts.ForceWithFlagA,
		subDir:        opts.SubDir,
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
	if ctx.removeLegacyTmp && ctx.dumpDir == "" {
		os.RemoveAll(ctx.genRoot)
	}
}

func (ctx *generateContext) Close() error {
	var closeErr error
	ctx.closeOnce.Do(func() {
		ctx.lifecycleMu.Lock()
		defer ctx.lifecycleMu.Unlock()
		ctx.closed = true
		// Safety if writeCases returned early before finishGenOrphans.
		ctx.releaseGenWrite()
		// Kind B expose lives in the product tree for go tool cover. GenerateOnly
		// leaves it for RunWorkspace (and leftover cleanup after all prepares).
		// The list is per gen root and shared across parallel PrepareTree — do
		// not strip here on generate failure or a sibling tree loses its files.
		if !ctx.generateOnly {
			closeErr = core.CleanupKindBMaterialized(ctx.genRoot)
		}
		ctx.removeTempsLocked()
	})
	return closeErr
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
	// Ephemeral build temps only. Kind B product files are covered by the
	// process handler in core (armed while materialized lists exist).
	if !ctx.removeLegacyTmp {
		return
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		// Hold lifecycleMu until process exit so concurrent writeCases cannot
		// MkdirAll/WriteFile a removed temp back into existence.
		ctx.lifecycleMu.Lock()
		ctx.closeOnce.Do(func() {
			ctx.closed = true
			ctx.reportKindBCleanup(core.CleanupAllKindBMaterialized())
			ctx.removeTempsLocked()
		})
		// Re-remove even if Close already ran without this lock held for exit.
		ctx.reportKindBCleanup(core.CleanupAllKindBMaterialized())
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
// Ephemeral build temps are written under lifecycleMu so SIGINT cleanup
// can RemoveAll without racing recreating writers.
//
// Always hierarchical unified (layout A): __droot, intermediates, leaf.go,
// suite/__allleaves/__registry. Parent/product internal → Kind B
// __doctest_internal_expose rewrite via ApplyInternalShimsAfterUnifiedGen.
func (ctx *generateContext) writeCases(cases []core.TreeCase, compileOnly bool) error {
	pkgName := "testcase"
	srcDir, origPkg, hasPkgUnderTest := core.ResolvePkgUnderTest(ctx.absRoot)
	if hasPkgUnderTest {
		pkgName = origPkg + "_tc"
	}

	if err := ctx.withGenLock(func() error {
		bridges, err := core.WriteGoModWithVendorBridges(ctx.genRoot, ctx.absModRoot, ctx.modPath, ctx.hasMod, ctx.assertImport, ctx.assertCacheDir, ctx.sessionImport, ctx.sessionCacheDir)
		if err != nil {
			return err
		}
		ctx.vendorBridges = bridges
		if ctx.verbose && ctx.w != nil {
			fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(filepath.Join(ctx.genRoot, "go.mod")))
		}
		// Framework packages (e.g. faas handlers) init via ProjectBasePath,
		// which requires src/ + config/ under the go test cwd (genRoot).
		if err := core.EnsureProjectBaseSymlinks(ctx.genRoot, ctx.absModRoot); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := ctx.withGenLock(func() error {
		if err := ctx.writeUnifiedCases(cases, compileOnly, pkgName, hasPkgUnderTest, srcDir, origPkg); err != nil {
			return err
		}
		if err := core.FlushGenManifest(ctx.genRoot); err != nil {
			return err
		}
		return ctx.finishGenOrphans()
	}); err != nil {
		return err
	}
	return nil
}

func (ctx *generateContext) reportKindBCleanup(err error) {
	if err == nil {
		return
	}
	w := ctx.w
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, "doctest: kind B cleanup: %v\n", err)
}

// joinKindBCleanupErr appends a Kind B product-file cleanup failure to a run
// error so a PASS cannot hide leftover __doctest_internal_expose files.
func joinKindBCleanupErr(err, cErr error) error {
	if cErr == nil {
		return err
	}
	if err != nil {
		return fmt.Errorf("%w; kind B cleanup: %v", err, cErr)
	}
	return fmt.Errorf("kind B cleanup: %w", cErr)
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
	if ctx.generateOnly || ctx.genBatch == nil {
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

	// Kind A: blank-import shims for leaf paths containing /internal/.
	// allleaves uses shim imports; real leaf packages stay at original paths.
	realLeafImports := append([]string(nil), leafImports...)
	leafImports = core.RewriteKindALeafImports(leafImports)

	// Suite/registry/allleaves under suiteRel (path-local when path-scoped).
	if err := core.WriteUnifiedTreeExtras(ctx.genRoot, suiteRel, ctx.absRoot, leafImports); err != nil {
		return err
	}
	// Emit kind A/B shim bodies + merge into vendor-gomod-overlay.json; rewrite
	// product-internal imports in gen sources to expose facades (kind B).
	if err := core.ApplyInternalShimsAfterUnifiedGen(ctx.genRoot, suiteRel, realLeafImports, cases, ctx.absModRoot, ctx.modPath); err != nil {
		return fmt.Errorf("internal shims: %w", err)
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

	if ctx.hasMod {
		if err := core.CondTidyGoMod(ctx.genRoot, ctx.goCache); err != nil {
			return err
		}
	}
	return core.FlushGenManifest(ctx.genRoot)
}

func (ctx *generateContext) scopedMultiRunDir(absRoot string) string {
	runDir := ctx.genRoot
	relRoot, err := filepath.Rel(ctx.absModRoot, absRoot)
	if err == nil && relRoot != "." {
		runDir = filepath.Join(ctx.genRoot, relRoot)
	}
	return runDir
}

// pathScopedGoTestPattern returns the go test package pattern for a gen-relative
// path scope (e.g. tree/mid → ./tree/mid/...). Suite packages still live under
// that prefix; ... selects them without hard-coding */suite.
func pathScopedGoTestPattern(suiteRel string) string {
	suiteRel = filepath.ToSlash(filepath.Clean(suiteRel))
	if suiteRel == "" || suiteRel == "." {
		return "./..."
	}
	return "./" + strings.TrimPrefix(suiteRel, "./") + "/..."
}

// packageArgsForCases returns go test package args for the unified suite layout.
// Full tree: single suite package. Path-scoped: ./<suiteRel>/...
func (ctx *generateContext) packageArgsForCases(runDir, absRoot string, cases []core.TreeCase) ([]string, error) {
	_ = absRoot
	_ = cases
	// Path-scoped: go test ./<suiteRel>/... (not a hard-coded */suite package).
	if ctx.isPathScoped() {
		return []string{pathScopedGoTestPattern(ctx.suiteRel())}, nil
	}
	// Full tree: single suite package under treeRel.
	suiteDir := core.UnifiedSuiteDirForTree(ctx.genRoot, ctx.treeRel())
	rel, err := filepath.Rel(runDir, suiteDir)
	if err != nil {
		return nil, fmt.Errorf("package path for suite: %w", err)
	}
	if rel == "." {
		return []string{"."}, nil
	}
	return []string{"./" + filepath.ToSlash(rel)}, nil
}
