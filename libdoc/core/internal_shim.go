// Package core internal shims: Kind A (scenario-path /internal/) and Kind B
// (product/parent module …/internal/… → __doctest_internal_expose facade).
//
// Production gen is always hierarchical unified (layout A). Parent/product
// internal imports never force classic multi-leaf .doctest_run_* compile; they
// are rewritten to Kind B expose packages via ApplyInternalShimsAfterUnifiedGen
// and merged into vendor-gomod-overlay.json.
//
// Kind B expose bodies live in __doctest_shim_store and are also materialized
// under the product module path (__doctest_internal_expose/…/expose.go) so
// go tool cover can open the logical path (overlay alone is not enough for
// cover). Paths are recorded in KindBMaterializedList and removed by
// CleanupKindBMaterialized at session end (RunWorkspace, leftover prepare
// roots, or generateContext.Close after a non-GenerateOnly run). Overlay
// remains for tools that prefer Replace maps. Overlay-only dirs still need
// -vet=off (NeedVetOff / InternalShimVetOffMarker) when cleanup has not run yet.
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Reserved path segments for virtual internal bridges (overlay materialization).
const (
	// Kind A: bridge under parent of scenario-path "internal" inside module testcase.
	DoctestInternalShimDir = "__doctest_internal_shim"
	// Kind B: facade under product module path (not named "internal").
	DoctestInternalExposeDir = "__doctest_internal_expose"
	// On-disk body store under gen root (never product VCS tree).
	InternalShimStoreDir = "__doctest_shim_store"
	// KindBMaterializedList under gen root lists product-tree expose.go paths
	// written for cover/compile; CleanupKindBMaterialized removes them.
	KindBMaterializedList = "doctest-kind-b-materialized"
)

// KindAShimImport maps a real gen import that contains /internal/ to a shim
// import under parent/__doctest_internal_shim/<suffix>. ok is false when no
// internal segment is present.
//
//	testcase/t/http/internal/post-succeeds
//	→ testcase/t/http/__doctest_internal_shim/post-succeeds
func KindAShimImport(realImport string) (shimImport string, ok bool) {
	parts := strings.Split(realImport, "/")
	for i, p := range parts {
		if p != "internal" {
			continue
		}
		// Collision: auto-suffix if user already used reserved name later — handled at emit.
		newParts := make([]string, 0, len(parts))
		newParts = append(newParts, parts[:i]...)
		newParts = append(newParts, DoctestInternalShimDir)
		newParts = append(newParts, parts[i+1:]...)
		return strings.Join(newParts, "/"), true
	}
	return "", false
}

// RewriteKindALeafImports returns allleaves import paths with kind A shims applied.
func RewriteKindALeafImports(leafImports []string) []string {
	out := make([]string, len(leafImports))
	for i, imp := range leafImports {
		if shim, ok := KindAShimImport(imp); ok {
			out[i] = shim
		} else {
			out[i] = imp
		}
	}
	return out
}

// ApplyInternalShimsAfterUnifiedGen emits kind A blank-import shims and kind B
// product expose facades, rewrites allleaves (already written with shim imports
// by caller for A), rewrites gen sources for B, and merges overlay Replace into
// genDir/vendor-gomod-overlay.json so tidy/test pick it up via existing flags.
//
// leafImports are the *original* real leaf imports (before kind A rewrite) so we
// can emit shim bodies. allLeavesDir is the directory containing all.go (already
// written with rewritten imports).
func ApplyInternalShimsAfterUnifiedGen(genRoot, suiteRel string, realLeafImports []string, cases []TreeCase, absModRoot, modPath string) error {
	if genRoot == "" {
		return nil
	}
	replace := make(map[string]string)

	// --- Kind A ---
	for _, realImp := range realLeafImports {
		shimImp, ok := KindAShimImport(realImp)
		if !ok {
			continue
		}
		if err := emitKindAShim(genRoot, realImp, shimImp, replace); err != nil {
			return err
		}
	}

	// --- Kind B ---
	if err := applyKindBExposes(genRoot, suiteRel, cases, absModRoot, modPath, replace); err != nil {
		return err
	}

	if len(replace) == 0 {
		_ = os.Remove(filepath.Join(genRoot, InternalShimVetOffMarker))
		return nil
	}
	// Overlay-only packages (kind B under product mod root) break `go test`
	// default vet (chdir into non-existent dir). Mark gen root so runners add -vet=off.
	if err := os.WriteFile(filepath.Join(genRoot, InternalShimVetOffMarker), []byte("kind-b-overlay\n"), 0644); err != nil {
		return err
	}
	return MergeReplaceIntoVendorGomodOverlay(genRoot, replace)
}

// InternalShimVetOffMarker is written when kind B overlay packages are active.
const InternalShimVetOffMarker = "doctest-internal-shim-vet-off"

// NeedVetOff reports whether go/xgo test should pass -vet=off for this gen root.
func NeedVetOff(genDir string) bool {
	if genDir == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(genDir, InternalShimVetOffMarker))
	return err == nil
}

func emitKindAShim(genRoot, realImp, shimImp string, replace map[string]string) error {
	// Kind A lives entirely under genRoot (module testcase). Write the shim
	// package on disk so go test/vet can chdir into it — pure overlay-only
	// packages fail with: vet: chdir … no such file or directory.
	_ = replace
	rel := strings.TrimPrefix(shimImp, "testcase/")
	if rel == shimImp {
		rel = shimImp
	}
	dir := filepath.Join(genRoot, filepath.FromSlash(rel))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	body := fmt.Sprintf("package shim\n\nimport _ %q\n", realImp)
	path := filepath.Join(dir, "shim.go")
	if prev, err := os.ReadFile(path); err == nil && string(prev) == body {
		return nil
	}
	return os.WriteFile(path, []byte(body), 0644)
}

func writeShimStoreBody(genRoot, kind, key, body string) (absPath string, err error) {
	sum := sha256.Sum256([]byte(kind + "\n" + key + "\n" + body))
	dir := filepath.Join(genRoot, InternalShimStoreDir, kind, hex.EncodeToString(sum[:12]))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "body.go")
	if prev, err := os.ReadFile(path); err == nil && string(prev) == body {
		return absPathOr(path), nil
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		return "", err
	}
	return absPathOr(path), nil
}

func absPathOr(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return a
}

// addOverlayKeyVariants registers Replace keys for p→body using darwin /private variants.
func addOverlayKeyVariants(replace map[string]string, virtFile, bodyPath string) {
	bodyAbs := absPathOr(bodyPath)
	// Prefer a body path that exists.
	if st, err := os.Stat(bodyAbs); err != nil || st.IsDir() {
		for _, b := range darwinPathVariants(bodyAbs) {
			if st, err := os.Stat(b); err == nil && !st.IsDir() {
				bodyAbs = b
				break
			}
		}
	}
	for _, k := range darwinPathVariants(absPathOr(virtFile)) {
		replace[k] = bodyAbs
	}
}

func darwinPathVariants(p string) []string {
	p = filepath.Clean(p)
	seen := map[string]bool{}
	var out []string
	add := func(s string) {
		if s == "" || seen[s] {
			return
		}
		seen[s] = true
		out = append(out, s)
	}
	add(p)
	if strings.HasPrefix(p, "/private/") {
		add(strings.TrimPrefix(p, "/private"))
	} else if strings.HasPrefix(p, "/var/") || strings.HasPrefix(p, "/tmp/") {
		add("/private" + p)
	}
	return out
}

// productInternalImport reports whether imp is someModule/internal/... eligible
// for Kind B expose facades. Parent-module internals are included (Kind B);
// only gen-module testcase paths are excluded (Kind A).
// parentModPath is retained for call-site compatibility but does not skip.
func productInternalImport(imp, parentModPath string) (productMod, internalPkg string, ok bool) {
	_ = parentModPath
	const marker = "/internal/"
	idx := strings.Index(imp, marker)
	if idx < 0 {
		return "", "", false
	}
	productMod = imp[:idx]
	if productMod == "" {
		return "", "", false
	}
	// Gen module testcase paths are kind A (scenario-path internal), not product.
	if productMod == "testcase" || strings.HasPrefix(imp, "testcase/") {
		return "", "", false
	}
	return productMod, imp, true
}

func applyKindBExposes(genRoot, suiteRel string, cases []TreeCase, absModRoot, parentModPath string, replace map[string]string) error {
	// Collect unique product internal imports from cases.
	needed := map[string]bool{} // full internal import path
	for _, tc := range cases {
		imports := collectImports(tc.SetupFiles, tc.AssertFile.GoBlock)
		for _, spec := range imports {
			if _, full, ok := productInternalImport(spec.Path, parentModPath); ok {
				needed[full] = true
			}
		}
	}
	if len(needed) == 0 {
		return nil
	}

	internals := make([]string, 0, len(needed))
	for p := range needed {
		internals = append(internals, p)
	}
	sort.Strings(internals)

	// Map internal import → expose import for rewrite.
	rewrites := map[string]string{}
	for _, internalImp := range internals {
		productMod, _, _ := productInternalImport(internalImp, parentModPath)
		exposeImp, err := emitKindBExpose(genRoot, absModRoot, productMod, internalImp, replace)
		if err != nil {
			return err
		}
		rewrites[internalImp] = exposeImp
	}

	// Rewrite import paths in generated Go under the suite tree.
	treeDir := genRoot
	if suiteRel != "" && suiteRel != "." {
		treeDir = filepath.Join(genRoot, suiteRel)
	}
	return rewriteImportsInTree(treeDir, rewrites)
}

func emitKindBExpose(genRoot, absModRoot, productMod, internalImp string, replace map[string]string) (exposeImp string, err error) {
	// expose path: productMod/__doctest_internal_expose/<tail after internal/>
	const marker = "/internal/"
	idx := strings.Index(internalImp, marker)
	if idx < 0 {
		return "", fmt.Errorf("internal import missing marker: %s", internalImp)
	}
	tail := internalImp[idx+len(marker):] // e.g. greet or foo/bar
	exposeImp = productMod + "/" + DoctestInternalExposeDir + "/" + tail

	productDir, err := resolveModuleDir(absModRoot, productMod)
	if err != nil {
		return "", fmt.Errorf("resolve product module %s: %w", productMod, err)
	}
	// Source dir of internal package on disk.
	internalRel := strings.TrimPrefix(internalImp, productMod+"/")
	internalDir := filepath.Join(productDir, filepath.FromSlash(internalRel))

	body, err := generateExposeSource(internalImp, internalDir)
	if err != nil {
		return "", fmt.Errorf("generate expose for %s: %w", internalImp, err)
	}
	bodyPath, err := writeShimStoreBody(genRoot, "b", internalImp, body)
	if err != nil {
		return "", err
	}

	// Logical path under product module root (package import path must resolve here).
	virtRel := DoctestInternalExposeDir + "/" + tail + "/expose.go"
	virtFile := filepath.Join(productDir, filepath.FromSlash(virtRel))
	// Materialize on disk: go tool cover opens this path without applying -overlay.
	// Keep shim-store + overlay as well for dual-path tooling.
	if err := os.MkdirAll(filepath.Dir(virtFile), 0755); err != nil {
		return "", fmt.Errorf("mkdir kind B expose %s: %w", virtFile, err)
	}
	if err := os.WriteFile(virtFile, []byte(body), 0644); err != nil {
		return "", fmt.Errorf("write kind B expose %s: %w", virtFile, err)
	}
	if err := recordKindBMaterialized(genRoot, virtFile); err != nil {
		return "", err
	}
	addOverlayKeyVariants(replace, virtFile, bodyPath)
	return exposeImp, nil
}

// recordKindBMaterialized appends an absolute product-tree expose path so
// CleanupKindBMaterialized can remove session-scoped files after go test/build.
func recordKindBMaterialized(genRoot, virtFile string) error {
	if genRoot == "" || virtFile == "" {
		return nil
	}
	abs, err := filepath.Abs(virtFile)
	if err != nil {
		abs = virtFile
	}
	if !isKindBMaterializedPath(abs) {
		return fmt.Errorf("record kind B materialized: not an expose path: %s", virtFile)
	}
	listPath := filepath.Join(genRoot, KindBMaterializedList)
	f, err := os.OpenFile(listPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("record kind B materialized: %w", err)
	}
	defer f.Close()
	if _, err := fmt.Fprintln(f, abs); err != nil {
		return fmt.Errorf("record kind B materialized: %w", err)
	}
	return nil
}

// CleanupKindBMaterialized removes product-module expose.go files listed under
// genRoot/KindBMaterializedList and prunes empty parent dirs up through
// __doctest_internal_expose. Safe when the list is missing (no Kind B).
func CleanupKindBMaterialized(genRoot string) error {
	if genRoot == "" {
		return nil
	}
	listPath := filepath.Join(genRoot, KindBMaterializedList)
	data, err := os.ReadFile(listPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var firstErr error
	var remaining []string
	for _, line := range strings.Split(string(data), "\n") {
		path := strings.TrimSpace(line)
		if path == "" {
			continue
		}
		if !isKindBMaterializedPath(path) {
			remaining = append(remaining, path)
			if firstErr == nil {
				firstErr = fmt.Errorf("cleanup kind B: refuse non-expose path %s", path)
			}
			continue
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			remaining = append(remaining, path)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		// Prune empty dirs up through __doctest_internal_expose only.
		dir := filepath.Dir(path)
		for dir != "" && dir != string(filepath.Separator) {
			base := filepath.Base(dir)
			entries, rdErr := os.ReadDir(dir)
			if rdErr != nil || len(entries) > 0 {
				break
			}
			if err := os.Remove(dir); err != nil {
				break
			}
			if base == DoctestInternalExposeDir {
				break
			}
			dir = filepath.Dir(dir)
		}
	}
	if len(remaining) > 0 {
		if err := os.WriteFile(listPath, []byte(strings.Join(remaining, "\n")+"\n"), 0644); err != nil && firstErr == nil {
			firstErr = err
		}
		return firstErr
	}
	_ = os.Remove(listPath)
	return firstErr
}

// isKindBMaterializedPath reports whether path is an absolute expose.go under
// __doctest_internal_expose (the only files emitKindBExpose writes).
func isKindBMaterializedPath(path string) bool {
	if path == "" || !filepath.IsAbs(path) {
		return false
	}
	if filepath.Base(path) != "expose.go" {
		return false
	}
	for _, p := range strings.Split(filepath.Clean(path), string(filepath.Separator)) {
		if p == DoctestInternalExposeDir {
			return true
		}
	}
	return false
}

func resolveModuleDir(workDir, modulePath string) (string, error) {
	if workDir == "" || modulePath == "" {
		return "", fmt.Errorf("workDir and modulePath required")
	}
	if dir, err := goListModuleDir(workDir, modulePath); err == nil && dir != "" {
		return dir, nil
	}
	// Fallback: go list fails when require was tidied away but replace remains
	// (common before gen imports the product package).
	if dir, err := resolveModuleDirFromGoMod(workDir, modulePath); err == nil && dir != "" {
		return dir, nil
	}
	return "", fmt.Errorf("resolve module %s under %s", modulePath, workDir)
}

func goListModuleDir(workDir, modulePath string) (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", modulePath)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go list -m %s: %w\n%s", modulePath, err, out)
	}
	dir := strings.TrimSpace(string(out))
	if dir == "" || dir == "|" {
		return "", fmt.Errorf("empty Dir for module %s", modulePath)
	}
	// go list may print warnings on stderr mixed — take last non-empty line
	lines := strings.Split(dir, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		l := strings.TrimSpace(lines[i])
		if l != "" && !strings.HasPrefix(l, "go:") {
			return l, nil
		}
	}
	return dir, nil
}

// resolveModuleDirFromGoMod reads replace/require path from go.mod.
func resolveModuleDirFromGoMod(workDir, modulePath string) (string, error) {
	data, err := os.ReadFile(filepath.Join(workDir, "go.mod"))
	if err != nil {
		return "", err
	}
	// replace example.com/app => ./app
	// replace example.com/app v0.0.0 => ../app
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "replace ") {
			continue
		}
		fields := strings.Fields(line)
		// replace PATH => DIR  OR  replace PATH VER => DIR
		if len(fields) < 4 || fields[1] != modulePath {
			continue
		}
		arrow := -1
		for i, f := range fields {
			if f == "=>" {
				arrow = i
				break
			}
		}
		if arrow < 0 || arrow+1 >= len(fields) {
			continue
		}
		rep := fields[arrow+1]
		if !filepath.IsAbs(rep) {
			rep = filepath.Join(workDir, rep)
		}
		abs, err := filepath.Abs(rep)
		if err != nil {
			return rep, nil
		}
		if st, err := os.Stat(abs); err != nil || !st.IsDir() {
			return "", fmt.Errorf("replace target not a dir: %s", abs)
		}
		return abs, nil
	}
	return "", fmt.Errorf("no replace for %s in go.mod", modulePath)
}

// exposeFunc is an exported func from one source file, with that file's import map
// (local alias → import path) so signatures can pull external packages.
type exposeFunc struct {
	fd        *ast.FuncDecl
	localPath map[string]string // import alias or default name → path
}

// generateExposeSource builds a facade package re-exporting exported funcs,
// types, vars, and consts from the product internal package so external modules
// can import the expose path. Package name matches the internal package.
//
// Func signatures that mention other product packages (e.g. model.Project) also
// emit those imports so the facade compiles under go test / go tool cover.
func generateExposeSource(internalImp, internalDir string) (string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, internalDir, func(fi os.FileInfo) bool {
		name := fi.Name()
		return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		return "", err
	}
	var pkg *ast.Package
	for name, p := range pkgs {
		if !strings.HasSuffix(name, "_test") {
			pkg = p
			break
		}
	}
	if pkg == nil {
		return "", fmt.Errorf("no package in %s", internalDir)
	}
	pkgName := pkg.Name
	srcAlias := "srcpkg"

	var types []string
	seenTy := map[string]bool{}
	var vars []string
	seenVar := map[string]bool{}
	var consts []string
	seenConst := map[string]bool{}
	var funcs []exposeFunc
	seenFn := map[string]bool{}
	// path → emit alias (stable across files; one import per path).
	pathToEmit := map[string]string{}
	reservedAlias := map[string]bool{srcAlias: true, pkgName: true}

	for _, f := range pkg.Files {
		localPath := fileImportAliasToPath(f, internalImp, internalDir)
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Recv != nil || d.Name == nil || !d.Name.IsExported() {
					continue
				}
				if seenFn[d.Name.Name] {
					continue
				}
				seenFn[d.Name.Name] = true
				funcs = append(funcs, exposeFunc{fd: d, localPath: localPath})
				// Register external packages used in this signature.
				for alias := range collectPkgSelectorsInFunc(d) {
					path, ok := localPath[alias]
					if !ok || path == "" || path == internalImp {
						continue
					}
					if _, exists := pathToEmit[path]; exists {
						continue
					}
					emit := importIdent(path, internalImp, internalDir)
					if reservedAlias[emit] || emitAliasTaken(pathToEmit, emit) {
						emit = uniqueEmitAlias(path, reservedAlias, pathToEmit)
					}
					pathToEmit[path] = emit
					reservedAlias[emit] = true
				}
			case *ast.GenDecl:
				switch d.Tok {
				case token.TYPE:
					for _, sp := range d.Specs {
						ts, ok := sp.(*ast.TypeSpec)
						if !ok || ts.Name == nil || !ts.Name.IsExported() {
							continue
						}
						if seenTy[ts.Name.Name] {
							continue
						}
						seenTy[ts.Name.Name] = true
						types = append(types, ts.Name.Name)
					}
				case token.VAR:
					for _, sp := range d.Specs {
						vs, ok := sp.(*ast.ValueSpec)
						if !ok {
							continue
						}
						for _, n := range vs.Names {
							if n == nil || !n.IsExported() || seenVar[n.Name] {
								continue
							}
							seenVar[n.Name] = true
							vars = append(vars, n.Name)
						}
					}
				case token.CONST:
					for _, sp := range d.Specs {
						vs, ok := sp.(*ast.ValueSpec)
						if !ok {
							continue
						}
						for _, n := range vs.Names {
							if n == nil || !n.IsExported() || seenConst[n.Name] {
								continue
							}
							seenConst[n.Name] = true
							consts = append(consts, n.Name)
						}
					}
				}
			}
		}
	}
	sort.Strings(types)
	sort.Strings(vars)
	sort.Strings(consts)
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].fd.Name.Name < funcs[j].fd.Name.Name
	})

	var b strings.Builder
	b.WriteString("// Code generated by doctest internal expose; DO NOT EDIT.\n")
	b.WriteString("package ")
	b.WriteString(pkgName)
	b.WriteString("\n\n")
	writeExposeImports(&b, srcAlias, internalImp, pathToEmit)
	b.WriteString("\n")
	for _, ty := range types {
		b.WriteString("type ")
		b.WriteString(ty)
		b.WriteString(" = ")
		b.WriteString(srcAlias)
		b.WriteString(".")
		b.WriteString(ty)
		b.WriteString("\n")
	}
	if len(types) > 0 {
		b.WriteString("\n")
	}
	for _, name := range vars {
		b.WriteString("var ")
		b.WriteString(name)
		b.WriteString(" = ")
		b.WriteString(srcAlias)
		b.WriteString(".")
		b.WriteString(name)
		b.WriteString("\n")
	}
	if len(vars) > 0 {
		b.WriteString("\n")
	}
	for _, name := range consts {
		// const X = srcpkg.X is valid when X is a constant in the source package.
		b.WriteString("const ")
		b.WriteString(name)
		b.WriteString(" = ")
		b.WriteString(srcAlias)
		b.WriteString(".")
		b.WriteString(name)
		b.WriteString("\n")
	}
	if len(consts) > 0 {
		b.WriteString("\n")
	}
	for _, ef := range funcs {
		fd := ef.fd
		sig := fd.Type
		// Map this file's import aliases → emit aliases for signature printing.
		localToEmit := map[string]string{}
		for local, path := range ef.localPath {
			if emit, ok := pathToEmit[path]; ok {
				localToEmit[local] = emit
			}
		}
		b.WriteString("func ")
		b.WriteString(fd.Name.Name)
		b.WriteString("(")
		b.WriteString(fieldListString(sig.Params, true, localToEmit))
		b.WriteString(")")
		if sig.Results != nil && len(sig.Results.List) > 0 {
			b.WriteString(" ")
			if len(sig.Results.List) == 1 && len(sig.Results.List[0].Names) == 0 {
				b.WriteString(exprString(sig.Results.List[0].Type, localToEmit))
			} else {
				b.WriteString("(")
				b.WriteString(fieldListString(sig.Results, false, localToEmit))
				b.WriteString(")")
			}
		}
		b.WriteString(" {\n\t")
		if sig.Results != nil && len(sig.Results.List) > 0 {
			b.WriteString("return ")
		}
		b.WriteString(srcAlias)
		b.WriteString(".")
		b.WriteString(fd.Name.Name)
		b.WriteString("(")
		b.WriteString(fieldListArgs(sig.Params))
		b.WriteString(")\n}\n\n")
	}
	return b.String(), nil
}

// fileImportAliasToPath maps import aliases (or resolved package names) to paths.
// Blank and dot imports are skipped (not used as signature qualifiers).
func fileImportAliasToPath(f *ast.File, internalImp, internalDir string) map[string]string {
	out := map[string]string{}
	if f == nil {
		return out
	}
	for _, imp := range f.Imports {
		if imp.Path == nil {
			continue
		}
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil || path == "" {
			continue
		}
		if imp.Name != nil {
			name := imp.Name.Name
			if name == "_" || name == "." {
				continue
			}
			out[name] = path
			continue
		}
		out[importIdent(path, internalImp, internalDir)] = path
	}
	return out
}

// importIdent is the identifier an unaliased import binds as: the package
// clause when the imported dir is resolvable next to internalDir, otherwise
// defaultImportAlias (last element, with /vN and name.vN heuristics).
func importIdent(importPath, internalImp, internalDir string) string {
	if dir := resolveImportDir(importPath, internalImp, internalDir); dir != "" {
		if name := readDirPackageName(dir); name != "" {
			return name
		}
	}
	return defaultImportAlias(importPath)
}

func resolveImportDir(importPath, internalImp, internalDir string) string {
	if importPath == "" || internalDir == "" {
		return ""
	}
	if importPath == internalImp {
		return internalDir
	}
	modRoot, modPath, ok := FindModuleRoot(internalDir)
	if !ok || modPath == "" {
		return ""
	}
	if importPath == modPath {
		return modRoot
	}
	if strings.HasPrefix(importPath, modPath+"/") {
		rel := filepath.FromSlash(strings.TrimPrefix(importPath, modPath+"/"))
		return filepath.Join(modRoot, rel)
	}
	return ""
}

func readDirPackageName(dir string) string {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		name := fi.Name()
		return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil || len(pkgs) == 0 {
		return ""
	}
	for name := range pkgs {
		if name != "" && name != "main" && !strings.HasSuffix(name, "_test") {
			return name
		}
	}
	return ""
}

// defaultImportAlias is the conventional unaliased name for an import path:
// last element, except /v2 → parent element and yaml.v3 → yaml.
func defaultImportAlias(path string) string {
	path = strings.TrimSuffix(path, "/")
	base := path
	if i := strings.LastIndex(path, "/"); i >= 0 {
		base = path[i+1:]
	}
	if isVersionElem(base) {
		parent := path
		if i := strings.LastIndex(path, "/"); i >= 0 {
			parent = path[:i]
			if j := strings.LastIndex(parent, "/"); j >= 0 {
				return parent[j+1:]
			}
			if parent != "" {
				return parent
			}
		}
	}
	if i := strings.LastIndex(base, "."); i > 0 && isVersionElem(base[i+1:]) {
		return base[:i]
	}
	return base
}

func isVersionElem(s string) bool {
	if len(s) < 2 || s[0] != 'v' {
		return false
	}
	for _, c := range s[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func emitAliasTaken(pathToEmit map[string]string, alias string) bool {
	for _, a := range pathToEmit {
		if a == alias {
			return true
		}
	}
	return false
}

func uniqueEmitAlias(path string, reserved map[string]bool, pathToEmit map[string]string) string {
	base := defaultImportAlias(path)
	// Prefer path-based uniqueness: replace / with _.
	cand := strings.ReplaceAll(strings.Trim(path, "/"), "/", "_")
	cand = strings.ReplaceAll(cand, ".", "_")
	cand = strings.ReplaceAll(cand, "-", "_")
	if cand == "" {
		cand = base + "_ext"
	}
	if !reserved[cand] && !emitAliasTaken(pathToEmit, cand) {
		return cand
	}
	for i := 2; i < 1000; i++ {
		a := fmt.Sprintf("%s_%d", base, i)
		if !reserved[a] && !emitAliasTaken(pathToEmit, a) {
			return a
		}
	}
	return base + "_ext"
}

// collectPkgSelectorsInFunc returns package identifiers used as X in X.Y type
// expressions within the function signature (params + results).
func collectPkgSelectorsInFunc(fd *ast.FuncDecl) map[string]bool {
	out := map[string]bool{}
	if fd == nil || fd.Type == nil {
		return out
	}
	collectPkgSelectorsFieldList(fd.Type.Params, out)
	collectPkgSelectorsFieldList(fd.Type.Results, out)
	return out
}

func collectPkgSelectorsFieldList(fl *ast.FieldList, out map[string]bool) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		collectPkgSelectorsExpr(f.Type, out)
	}
}

func collectPkgSelectorsExpr(e ast.Expr, out map[string]bool) {
	if e == nil {
		return
	}
	switch t := e.(type) {
	case *ast.SelectorExpr:
		if id, ok := t.X.(*ast.Ident); ok {
			out[id.Name] = true
		} else {
			collectPkgSelectorsExpr(t.X, out)
		}
	case *ast.StarExpr:
		collectPkgSelectorsExpr(t.X, out)
	case *ast.ArrayType:
		collectPkgSelectorsExpr(t.Len, out)
		collectPkgSelectorsExpr(t.Elt, out)
	case *ast.Ellipsis:
		collectPkgSelectorsExpr(t.Elt, out)
	case *ast.MapType:
		collectPkgSelectorsExpr(t.Key, out)
		collectPkgSelectorsExpr(t.Value, out)
	case *ast.ChanType:
		collectPkgSelectorsExpr(t.Value, out)
	case *ast.FuncType:
		collectPkgSelectorsFieldList(t.Params, out)
		collectPkgSelectorsFieldList(t.Results, out)
	case *ast.InterfaceType:
		// ignore method sets for expose wrappers
	case *ast.IndexExpr:
		collectPkgSelectorsExpr(t.X, out)
		collectPkgSelectorsExpr(t.Index, out)
	case *ast.IndexListExpr:
		collectPkgSelectorsExpr(t.X, out)
		for _, ix := range t.Indices {
			collectPkgSelectorsExpr(ix, out)
		}
	}
}

func writeExposeImports(b *strings.Builder, srcAlias, internalImp string, pathToEmit map[string]string) {
	// Sort external paths for stable output.
	paths := make([]string, 0, len(pathToEmit))
	for p := range pathToEmit {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	if len(paths) == 0 {
		b.WriteString("import ")
		b.WriteString(srcAlias)
		b.WriteString(" ")
		b.WriteString(strconv.Quote(internalImp))
		b.WriteString("\n")
		return
	}
	b.WriteString("import (\n\t")
	b.WriteString(srcAlias)
	b.WriteString(" ")
	b.WriteString(strconv.Quote(internalImp))
	b.WriteString("\n")
	for _, p := range paths {
		emit := pathToEmit[p]
		b.WriteString("\t")
		if emit == defaultImportAlias(p) {
			b.WriteString(strconv.Quote(p))
		} else {
			b.WriteString(emit)
			b.WriteString(" ")
			b.WriteString(strconv.Quote(p))
		}
		b.WriteString("\n")
	}
	b.WriteString(")\n")
}

func fieldListString(fl *ast.FieldList, withNames bool, localToEmit map[string]string) string {
	if fl == nil {
		return ""
	}
	var parts []string
	for _, f := range fl.List {
		ty := exprString(f.Type, localToEmit)
		if !withNames || len(f.Names) == 0 {
			parts = append(parts, ty)
			continue
		}
		for _, n := range f.Names {
			parts = append(parts, n.Name+" "+ty)
		}
	}
	return strings.Join(parts, ", ")
}

func fieldListArgs(fl *ast.FieldList) string {
	if fl == nil {
		return ""
	}
	var names []string
	for _, f := range fl.List {
		if len(f.Names) == 0 {
			names = append(names, "_")
			continue
		}
		for _, n := range f.Names {
			names = append(names, n.Name)
		}
	}
	return strings.Join(names, ", ")
}

// exprString prints a type expression. localToEmit rewrites package selectors
// (file-local import aliases → facade emit aliases).
func exprString(e ast.Expr, localToEmit map[string]string) string {
	if e == nil {
		return ""
	}
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprString(t.X, localToEmit)
	case *ast.SelectorExpr:
		if id, ok := t.X.(*ast.Ident); ok {
			if emit, ok := localToEmit[id.Name]; ok {
				return emit + "." + t.Sel.Name
			}
			return id.Name + "." + t.Sel.Name
		}
		return exprString(t.X, localToEmit) + "." + t.Sel.Name
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + exprString(t.Elt, localToEmit)
		}
		return "[" + exprString(t.Len, localToEmit) + "]" + exprString(t.Elt, localToEmit)
	case *ast.Ellipsis:
		return "..." + exprString(t.Elt, localToEmit)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.MapType:
		return "map[" + exprString(t.Key, localToEmit) + "]" + exprString(t.Value, localToEmit)
	case *ast.ChanType:
		return "chan " + exprString(t.Value, localToEmit)
	case *ast.FuncType:
		return "func"
	case *ast.BasicLit:
		return t.Value
	case *ast.IndexExpr:
		return exprString(t.X, localToEmit) + "[" + exprString(t.Index, localToEmit) + "]"
	case *ast.IndexListExpr:
		var parts []string
		for _, ix := range t.Indices {
			parts = append(parts, exprString(ix, localToEmit))
		}
		return exprString(t.X, localToEmit) + "[" + strings.Join(parts, ", ") + "]"
	default:
		return "interface{}"
	}
}

func rewriteImportsInTree(treeDir string, rewrites map[string]string) error {
	if len(rewrites) == 0 {
		return nil
	}
	return filepath.Walk(treeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// skip shim store
			if info.Name() == InternalShimStoreDir || info.Name() == vendorGomodOverlayDir {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := string(data)
		orig := s
		for old, neu := range rewrites {
			// Replace quoted import paths.
			s = strings.ReplaceAll(s, strconv.Quote(old), strconv.Quote(neu))
		}
		if s == orig {
			return nil
		}
		return os.WriteFile(path, []byte(s), info.Mode())
	})
}
