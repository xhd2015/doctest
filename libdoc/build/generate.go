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
	modPath         string
	hasMod          bool
	genRoot         string
	compileRoot     string
	dumpDir         string
	internalCompile bool
	removeLegacyTmp bool
	closeOnce       sync.Once
}

func newGenerateContext(dir string, opts core.Options, cases []core.TreeCase, w io.Writer, forBuild bool, verbose bool) (*generateContext, error) {
	absRoot, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	modRoot, modPath, hasMod := core.FindModuleRoot(absRoot)
	absModRoot, _ := core.MappingGenRoot(absRoot)

	ctx := &generateContext{
		w:       w,
		verbose: verbose,
		absRoot: absRoot,
		absModRoot: absModRoot,
		modPath: modPath,
		hasMod:  hasMod,
		dumpDir: opts.GenDir,
		internalCompile: hasMod && core.CasesImportInternalPackage(cases, modPath),
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
			genRoot, _, cacheErr = core.CacheMappingGenRoot(absRoot)
			if cacheErr != nil {
				return nil, cacheErr
			}
		}
	}
	if err := os.MkdirAll(genRoot, 0755); err != nil {
		return nil, err
	}
	ctx.genRoot = genRoot
	return ctx, nil
}

func (ctx *generateContext) Close() {
	ctx.closeOnce.Do(func() {
		if ctx.compileRoot != "" {
			os.RemoveAll(ctx.compileRoot)
		}
		if ctx.removeLegacyTmp && ctx.dumpDir == "" {
			os.RemoveAll(ctx.genRoot)
		}
	})
}

func (ctx *generateContext) installInterruptCleanup() {
	if ctx.compileRoot == "" && !ctx.removeLegacyTmp {
		return
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		ctx.Close()
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

func (ctx *generateContext) writeCases(cases []core.TreeCase, compileOnly bool) error {
	pkgName := "testcase"
	srcDir, origPkg, hasPkgUnderTest := core.ResolvePkgUnderTest(ctx.absRoot)
	if hasPkgUnderTest {
		pkgName = origPkg + "_tc"
	}

	if !ctx.internalCompile {
		if err := core.WriteGoMod(ctx.genRoot, ctx.absModRoot, ctx.modPath, ctx.hasMod); err != nil {
			return err
		}
		if ctx.verbose && ctx.w != nil {
			fmt.Fprintf(ctx.w, "→ %s\n", pathfmt.Short(filepath.Join(ctx.genRoot, "go.mod")))
		}
	}

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

		testPath, err := core.WriteGeneratedCase(leafDir, tc, compileOnly, pkgName, ctx.absRoot)
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
	}

	if ctx.hasMod && !ctx.internalCompile {
		if err := core.CondTidyGoMod(ctx.genRoot); err != nil {
			return err
		}
	}
	return nil
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