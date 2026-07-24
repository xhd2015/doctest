package build

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
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
	// unifiedMode: hierarchical ref packages + one suite package per DOCTEST tree.
	// Always true for normal generation; false only for internal-compile trees
	// (module-internal import path layout still uses classic AssembleTestSource).
	unifiedMode bool
	closeOnce   sync.Once
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
	// GenBatch: optional -a wipe-once; attach for emit-set orphan recording.
	if ctx.genBatch != nil {
		ctx.genBatch.Attach(genRoot)
		if ctx.forceA {
			if err := ctx.genBatch.WipeOnce(genRoot); err != nil {
				return nil, err
			}
		}
	} else if opts.ForceWithFlagA {
		// Library -a without shared GenBatch: still wipe this gen root.
		if err := core.WipeGenRoot(genRoot); err != nil {
			return nil, err
		}
	}
	ctx.genRoot = genRoot
	return ctx, nil
}

// removeTempsLocked deletes interrupt-scoped temps. Caller must hold lifecycleMu.
func (ctx *generateContext) removeTempsLocked() {
	if ctx.modfilePath != "" {
		os.Remove(ctx.modfilePath)
	}
	if ctx.compileRoot != "" {
		os.RemoveAll(ctx.compileRoot)
	}
	if ctx.removeLegacyTmp && ctx.dumpDir == "" {
		os.RemoveAll(ctx.genRoot)
	}
}

func (ctx *generateContext) Close() {
	ctx.closeOnce.Do(func() {
		ctx.lifecycleMu.Lock()
		defer ctx.lifecycleMu.Unlock()
		ctx.closed = true
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

// finishGenOrphans reconciles throwaway gen: delete paths not in this batch's
// emit set. Deferred when GenerateOnly (multi-tree workspace prunes later).
func (ctx *generateContext) finishGenOrphans() error {
	if ctx.internalCompile || ctx.generateOnly || ctx.genBatch == nil {
		return nil
	}
	return ctx.genBatch.PruneAttached(ctx.genRoot)
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
	rootDir := core.RefRootDirForTree(ctx.genRoot, treeRel)
	rootImport := core.RefRootImportForTree(treeRel)
	registryImport := core.UnifiedRegistryImportForTree(treeRel)
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

	if err := core.WriteUnifiedTreeExtras(ctx.genRoot, treeRel, ctx.absRoot, leafImports); err != nil {
		return err
	}
	// Workspace registration plane: each tree's __wreg registers into
	// __workspace/__registry (fan-in rewritten when multi-root runs).
	if err := core.WriteTreeWreg(ctx.genRoot, treeRel, ctx.absRoot); err != nil {
		return fmt.Errorf("write tree wreg: %w", err)
	}
	if ctx.verbose && ctx.w != nil {
		fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(core.UnifiedRegistryDirForTree(ctx.genRoot, treeRel)))
		fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(core.UnifiedAllLeavesDirForTree(ctx.genRoot, treeRel)))
		fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(core.UnifiedSuiteDirForTree(ctx.genRoot, treeRel)))
		fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(core.TreeWregDirForTree(ctx.genRoot, treeRel)))
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
		// Single suite package → one test binary per DOCTEST tree.
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
