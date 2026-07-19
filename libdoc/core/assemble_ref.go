package core

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/imports"
)

// Ref-mode package layout under gen root (module path is always "testcase"
// when WriteGoMod is used for the external gen tree):
//
//	genRoot/
//	  go.mod
//	  <treeRel>/__droot/droot.go   package droot — types, Run, root helpers
//	  <treeRel>/<leaf>/…_test.go   thin tests importing testcase/<treeRel>/__droot
//
// treeRel is the path of the doctest root relative to the module root (or "." →
// __droot at gen root). This keeps multi-tree ./... + shared GenDir (cold-cache)
// from overwriting each other's __droot packages.
const (
	RefRootDirName = "__droot"
	RefRootPkgName = "droot"
	// RefRootImportPath is the legacy flat import used when treeRel is ".".
	RefRootImportPath = "testcase/__droot"
)

// RefRootImportForTree returns the go import path for the tree-scoped droot package.
// treeRel is filepath.Rel(modRoot, doctestRoot); empty/"." → testcase/__droot.
func RefRootImportForTree(treeRel string) string {
	treeRel = filepath.ToSlash(filepath.Clean(treeRel))
	if treeRel == "" || treeRel == "." {
		return RefRootImportPath
	}
	return "testcase/" + strings.TrimPrefix(treeRel, "./") + "/" + RefRootDirName
}

// RefRootDirForTree returns the filesystem directory for __droot under genRoot.
func RefRootDirForTree(genRoot, treeRel string) string {
	treeRel = filepath.Clean(treeRel)
	if treeRel == "" || treeRel == "." {
		return filepath.Join(genRoot, RefRootDirName)
	}
	return filepath.Join(genRoot, treeRel, RefRootDirName)
}

// isRootSetupDoc reports whether a setup-chain document belongs to the tree
// root (shared package) rather than a leaf/intermediate directory.
func isRootSetupDoc(doc SetupDocument) bool {
	p := filepath.ToSlash(doc.Path)
	return p == "DOCTEST.md" || p == "SETUP.md" || p == ""
}

// SplitRefSetupDocs partitions the setup chain into root vs non-root docs.
func SplitRefSetupDocs(setupFiles []SetupDocument) (rootDocs, leafDocs []SetupDocument) {
	for _, doc := range setupFiles {
		if isRootSetupDoc(doc) {
			rootDocs = append(rootDocs, doc)
		} else {
			leafDocs = append(leafDocs, doc)
		}
	}
	return rootDocs, leafDocs
}

// collectRootSymbolRenames maps unexported root package symbols to exported
// names so leaf packages can reference them as droot.ExportedName.
// Exported symbols map to themselves. Types Request/Response are excluded
// (already qualified separately as droot.Request).
func collectRootSymbolRenames(rootDocs []SetupDocument) map[string]string {
	renames := make(map[string]string)
	add := func(name string) {
		if name == "" || name == "_" ||
			name == "Request" || name == "Response" || name == "Setup" || name == "Run" || name == "Assert" {
			return
		}
		// Never rewrite blank identifiers or tiny builtins via symbol map.
		if name == "err" || name == "ok" || name == "t" || name == "req" || name == "resp" {
			return
		}
		if token.IsExported(name) {
			renames[name] = name
			return
		}
		renames[name] = exportIdent(name)
	}
	for _, doc := range rootDocs {
		if doc.GoBlock == nil {
			continue
		}
		for _, h := range doc.GoBlock.Helpers {
			add(h.Name)
		}
		// Methods stay package-local on root types; still export name if leaf needed (rare).
		for _, m := range doc.GoBlock.Methods {
			// Method name only; receiver type stays in droot.
			add(m.Name)
		}
		for _, decl := range doc.GoBlock.VarDecls {
			for _, n := range declIdentNames(decl) {
				add(n)
			}
		}
		for _, decl := range doc.GoBlock.Consts {
			for _, n := range declIdentNames(decl) {
				add(n)
			}
		}
	}
	return renames
}

func exportIdent(name string) string {
	if name == "" || token.IsExported(name) {
		return name
	}
	r, size := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(r)) + name[size:]
}

func declIdentNames(decl string) []string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", "package p\n"+decl, 0)
	if err != nil {
		return nil
	}
	var names []string
	for _, d := range f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, id := range vs.Names {
				names = append(names, id.Name)
			}
		}
	}
	return names
}

// rewriteBareIdents renames bare identifiers according to renames (old→new).
// Uses go/scanner so string/rune/comment contents are left untouched (critical:
// fixture generators embed Go source in string literals).
func rewriteBareIdents(src string, renames map[string]string) string {
	if len(renames) == 0 || src == "" {
		return src
	}
	return scanReplaceIdents(src, func(name string, afterDot bool) string {
		if afterDot {
			return name
		}
		if neu, ok := renames[name]; ok {
			return neu
		}
		return name
	})
}

// qualifyRootSymbols prefixes bare root symbols with alias.ExportedName
// (e.g. runType → droot.RunType). Skips string/comment content.
func qualifyRootSymbols(src, alias string, renames map[string]string) string {
	if (len(renames) == 0 && alias == "") || src == "" {
		return src
	}
	// Map both bare and already-exported forms to the exported form.
	lookup := map[string]string{}
	for old, neu := range renames {
		lookup[old] = neu
		lookup[neu] = neu
	}
	return scanReplaceIdents(src, func(name string, afterDot bool) string {
		if afterDot || name == "_" || name == "" {
			return name
		}
		exported, ok := lookup[name]
		if !ok {
			return name
		}
		return alias + "." + exported
	})
}

// scanReplaceIdents rewrites IDENT tokens via repl, skipping strings/comments.
// afterDot is true when the previous non-comment token was '.'.
// Token spans are taken from consecutive scan offsets so operator tokens
// (empty lit) cannot desync the rewrite from the source.
func scanReplaceIdents(src string, repl func(name string, afterDot bool) string) string {
	type tok struct {
		off  int
		tok  token.Token
		lit  string
	}
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	var s scanner.Scanner
	s.Init(file, []byte(src), nil, scanner.ScanComments)
	var toks []tok
	for {
		pos, t, lit := s.Scan()
		if t == token.EOF {
			break
		}
		// Skip automatically inserted semicolons — keep original newlines via gaps.
		if t == token.SEMICOLON && lit == "\n" {
			continue
		}
		toks = append(toks, tok{off: fset.File(pos).Offset(pos), tok: t, lit: lit})
	}
	var out strings.Builder
	out.Grow(len(src) + 64)
	prevDot := false
	lastEnd := 0
	for i, tk := range toks {
		end := len(src)
		if i+1 < len(toks) {
			end = toks[i+1].off
		}
		// Prefer source slice for non-IDENT so spacing stays exact.
		if tk.off > lastEnd {
			out.WriteString(src[lastEnd:tk.off])
		}
		switch tk.tok {
		case token.IDENT:
			// Token text length is len(lit); do not use end-off (includes trailing space).
			neu := repl(tk.lit, prevDot)
			out.WriteString(neu)
			lastEnd = tk.off + len(tk.lit)
			prevDot = false
		case token.PERIOD:
			out.WriteByte('.')
			lastEnd = tk.off + 1
			prevDot = true
		default:
			// Copy exact source bytes for this token (and only the token, not trailing space).
			tokEnd := tk.off
			if tk.lit != "" {
				tokEnd = tk.off + len(tk.lit)
			} else {
				// Operators: length of token.String(), e.g. ":=", "..."
				tokEnd = tk.off + len(tk.tok.String())
				if tokEnd > end {
					tokEnd = end
				}
				// Prefer raw source slice when it matches.
				if tokEnd <= len(src) {
					out.WriteString(src[tk.off:tokEnd])
					lastEnd = tokEnd
					prevDot = false
					continue
				}
			}
			if tokEnd > len(src) {
				tokEnd = len(src)
			}
			out.WriteString(src[tk.off:tokEnd])
			lastEnd = tokEnd
			prevDot = false
		}
	}
	if lastEnd < len(src) {
		out.WriteString(src[lastEnd:])
	}
	return out.String()
}

// AssembleRefRootSource emits the shared root package: types, consts/vars,
// helpers/methods, and package-level Run from root setup documents.
// Unexported package-level symbols are exported so leaf packages can use them.
func AssembleRefRootSource(rootDocs []SetupDocument, pkgName string) (string, error) {
	if pkgName == "" {
		pkgName = RefRootPkgName
	}
	var run *FuncSnippet
	for _, doc := range rootDocs {
		if doc.GoBlock != nil && doc.GoBlock.Run != nil {
			runCopy := *doc.GoBlock.Run
			run = &runCopy
			break
		}
	}
	if run == nil {
		return "", fmt.Errorf("missing Run(t *testing.T, req *Request) (*Response, error) in root setup chain")
	}

	renames := collectRootSymbolRenames(rootDocs)

	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteString("\n\n")

	importsMap := collectImports(rootDocs, GoBlock{})
	if _, ok := importsMap["testing"]; !ok {
		importsMap["testing"] = &ImportSpec{Path: "testing"}
	}
	if _, ok := importsMap[sessionImportPath]; !ok {
		importsMap[sessionImportPath] = &ImportSpec{Path: sessionImportPath}
	}
	writeImportBlock(&buf, importsMap)

	// Session context is injected as d *session.Doctest on Setup/Run — no package free DOCTEST_* vars.

	// Emit types/methods/vars/helpers then rewrite unexported symbols to exported.
	var body strings.Builder
	writePackageLevelTypesAndMethods(&body, rootDocs, GoBlock{})
	writePackageLevelConstVars(&body, rootDocs, GoBlock{})
	writePackageLevelHelpers(&body, rootDocs, GoBlock{})

	for i, doc := range rootDocs {
		if doc.GoBlock == nil || doc.GoBlock.Setup == nil {
			continue
		}
		fn := *doc.GoBlock.Setup
		if fn.Name == "" || fn.Name == "Setup" {
			fn.Name = fmt.Sprintf("RootSetup%d", i)
		}
		fn.Params = ensureDoctestParam(fn.Params)
		writePackageFunc(&body, fn)
		body.WriteString("\n")
	}

	run.Name = "Run"
	run.Params = ensureDoctestParam(run.Params)
	writePackageFunc(&body, *run)
	body.WriteString("\n")

	// Export unexported package-level symbols (vars/helpers) and rewrite refs.
	bodyStr := rewriteBareIdents(body.String(), renames)
	buf.WriteString(bodyStr)
	return buf.String(), nil
}

// AssembleRefLeafTestSource emits a thin leaf *_test.go that imports the root
// package. Leaf must not redefine root types or root helpers.
func AssembleRefLeafTestSource(tc TreeCase, compileOnly bool, pkgName, docTestRoot, rootImport, rootAlias string) (string, error) {
	if pkgName == "" {
		pkgName = "testcase"
	}
	if rootImport == "" {
		rootImport = RefRootImportPath
	}
	if rootAlias == "" {
		rootAlias = RefRootPkgName
	}

	rootDocs, leafDocs := SplitRefSetupDocs(tc.SetupFiles)
	renames := collectRootSymbolRenames(rootDocs)
	rootTypes := collectRootTypeNames(rootDocs)

	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteString("\n\n")

	importsMap := collectImports(leafDocs, tc.AssertFile.GoBlock)
	for _, pkg := range []string{"testing", "syscall", sessionImportPath} {
		if _, ok := importsMap[pkg]; !ok {
			importsMap[pkg] = &ImportSpec{Path: pkg}
		}
	}
	importsMap[rootAlias+"\x00"+rootImport] = &ImportSpec{Name: rootAlias, Path: rootImport}
	writeImportBlock(&buf, importsMap)

	// Leaf-only types/helpers — rewrite any references to root symbols.
	var leafBlob strings.Builder
	writePackageLevelTypesAndMethods(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	writePackageLevelConstVars(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	writePackageLevelHelpers(&leafBlob, leafDocs, tc.AssertFile.GoBlock)
	leafTop := qualifyRootSymbols(leafBlob.String(), rootAlias, renames)
	// Root-declared types in leaf helpers/types become droot.X
	leafTop = qualifyRootTypes(leafTop, rootAlias, rootTypes)
	buf.WriteString(leafTop)

	buf.WriteString("func ")
	buf.WriteString(TestFuncName(tc))
	buf.WriteString("(t *testing.T) {\n")

	writeDoctestDConstruct(&buf, docTestRoot, tc.Path)

	buf.WriteString("\tRun := ")
	buf.WriteString(rootAlias)
	buf.WriteString(".Run\n")
	buf.WriteString(fmt.Sprintf("\treq := &%s.Request{}\n", rootAlias))

	for i, doc := range rootDocs {
		if doc.GoBlock == nil || doc.GoBlock.Setup == nil {
			continue
		}
		name := fmt.Sprintf("RootSetup%d", i)
		buf.WriteString(fmt.Sprintf("\tif err := %s.%s(t, d, req); err != nil {\n", rootAlias, name))
		buf.WriteString(fmt.Sprintf("\t\tt.Fatalf(\"%s failed: %%v\", err)\n", escapeString(doc.Path)))
		buf.WriteString("\t}\n")
	}

	for i, doc := range leafDocs {
		if doc.GoBlock == nil || doc.GoBlock.Setup == nil {
			continue
		}
		name := fmt.Sprintf("setup%d", i)
		fn := *doc.GoBlock.Setup
		fn.Params = qualifyRootTypes(fn.Params, rootAlias, rootTypes)
		fn.Results = qualifyRootTypes(fn.Results, rootAlias, rootTypes)
		fn.ResultTypes = qualifyRootTypes(fn.ResultTypes, rootAlias, rootTypes)
		fn.ClosureResults = qualifyRootTypes(fn.ClosureResults, rootAlias, rootTypes)
		fn.Body = qualifyRootTypesInBody(fn.Body, rootAlias, rootTypes)
		fn.Body = qualifyRootSymbols(fn.Body, rootAlias, renames)
		fn.Params = qualifyRootSymbols(fn.Params, rootAlias, renames)
		fn.Params = ensureDoctestParam(fn.Params)
		writeFuncClosure(&buf, name, fn)
		buf.WriteString(fmt.Sprintf("\tif err := %s(t, d, req); err != nil {\n", name))
		buf.WriteString(fmt.Sprintf("\t\tt.Fatalf(\"%s failed: %%v\", err)\n", escapeString(doc.Path)))
		buf.WriteString("\t}\n")
	}

	assertFn := *tc.AssertFile.GoBlock.Assert
	assertFn.Params = qualifyRootTypes(assertFn.Params, rootAlias, rootTypes)
	assertFn.Results = qualifyRootTypes(assertFn.Results, rootAlias, rootTypes)
	assertFn.ResultTypes = qualifyRootTypes(assertFn.ResultTypes, rootAlias, rootTypes)
	assertFn.ClosureResults = qualifyRootTypes(assertFn.ClosureResults, rootAlias, rootTypes)
	assertFn.Body = qualifyRootTypesInBody(assertFn.Body, rootAlias, rootTypes)
	assertFn.Body = qualifyRootSymbols(assertFn.Body, rootAlias, renames)
	assertFn.Params = ensureDoctestParam(assertFn.Params)
	writeFuncClosure(&buf, "assert", assertFn)

	buf.WriteString("\t_ = Run\n")
	helperNames := collectHelperNames(leafDocs, tc.AssertFile.GoBlock)
	for _, name := range helperNames {
		buf.WriteString(fmt.Sprintf("\t_ = %s\n", name))
	}

	if compileOnly {
		buf.WriteString("\t// compileOnly\n")
		buf.WriteString("\t_ = d\n")
		buf.WriteString("\t_ = req\n")
		buf.WriteString("\t_ = assert\n")
		buf.WriteString(fmt.Sprintf("\tvar resp *%s.Response\n", rootAlias))
		buf.WriteString("\tvar runErr error\n")
		buf.WriteString("\t_ = resp\n")
		buf.WriteString("\t_ = runErr\n")
		buf.WriteString("}\n")
		return buf.String(), nil
	}

	buf.WriteString(fmt.Sprintf("\tresp, runErr := %s.Run(t, d, req)\n", rootAlias))
	buf.WriteString("\tassert(t, d, req, resp, runErr)\n")
	buf.WriteString("}\n")
	return buf.String(), nil
}

func writeImportBlock(buf *strings.Builder, imports map[string]*ImportSpec) {
	if len(imports) == 0 {
		return
	}
	buf.WriteString("import (\n")
	importList := make([]*ImportSpec, 0, len(imports))
	for _, spec := range imports {
		importList = append(importList, spec)
	}
	sort.Slice(importList, func(i, j int) bool {
		if importList[i].Path != importList[j].Path {
			return importList[i].Path < importList[j].Path
		}
		return importList[i].Name < importList[j].Name
	})
	for _, spec := range importList {
		if spec.Name != "" {
			buf.WriteString("\t" + spec.Name + " \"" + spec.Path + "\"\n")
		} else {
			buf.WriteString("\t\"" + spec.Path + "\"\n")
		}
	}
	buf.WriteString(")\n\n")
}

// collectRootTypeNames returns type names declared in root docs (Request, Response, WarnCase, …).
func collectRootTypeNames(rootDocs []SetupDocument) []string {
	seen := map[string]bool{}
	var names []string
	for _, doc := range rootDocs {
		if doc.GoBlock == nil {
			continue
		}
		for name := range doc.GoBlock.Types {
			if name == "" || seen[name] {
				continue
			}
			seen[name] = true
			names = append(names, name)
		}
	}
	sort.Slice(names, func(i, j int) bool { return len(names[i]) > len(names[j]) })
	return names
}

// qualifyRootTypes rewrites bare root type identifiers to rootAlias.X
// (Request, Response, and any other types declared in root docs).
// Skips string/comment content so fixture generators keep raw `*Request` text.
func qualifyRootTypes(s, rootAlias string, rootTypeNames []string) string {
	if s == "" {
		return s
	}
	if len(rootTypeNames) == 0 {
		rootTypeNames = []string{"Request", "Response"}
	}
	typeSet := map[string]bool{}
	for _, t := range rootTypeNames {
		typeSet[t] = true
	}
	return scanReplaceIdents(s, func(name string, afterDot bool) string {
		if afterDot || !typeSet[name] {
			return name
		}
		return rootAlias + "." + name
	})
}

func qualifyRootTypesInBody(body, rootAlias string, rootTypeNames []string) string {
	return qualifyRootTypes(body, rootAlias, rootTypeNames)
}

// WriteRefTree writes the shared root package once and thin leaf tests for each case.
func WriteRefTree(genRoot string, cases []TreeCase, docTestRoot string, compileOnly bool, pkgName string) error {
	if len(cases) == 0 {
		return fmt.Errorf("WriteRefTree: no cases")
	}
	if pkgName == "" {
		pkgName = "testcase"
	}

	rootDocs, _ := SplitRefSetupDocs(cases[0].SetupFiles)
	if len(rootDocs) == 0 {
		rootDocs = cases[0].SetupFiles
	}
	rootSrc, err := AssembleRefRootSource(rootDocs, RefRootPkgName)
	if err != nil {
		return err
	}
	// Tree-relative layout: put __droot next to leaves under tree path.
	// docTestRoot is the tree root; without mod root we use flat __droot.
	rootDir := filepath.Join(genRoot, RefRootDirName)
	rootImport := RefRootImportPath
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return err
	}
	rootPath := filepath.Join(rootDir, "droot.go")
	if err := WriteFormattedGo(rootPath, rootSrc); err != nil {
		return fmt.Errorf("write ref root package: %w", err)
	}

	for _, tc := range cases {
		leafDir := genRoot
		if tc.Path != "" {
			leafDir = filepath.Join(genRoot, tc.Path)
		}
		if _, err := WriteRefLeafCase(leafDir, tc, compileOnly, pkgName, docTestRoot, rootImport); err != nil {
			return err
		}
	}
	return nil
}

// WriteRefLeafCase writes a thin leaf test file under leafDir.
func WriteRefLeafCase(leafDir string, tc TreeCase, compileOnly bool, pkgName, docTestRoot, rootImport string) (string, error) {
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		return "", err
	}
	if rootImport == "" {
		rootImport = RefRootImportPath
	}
	src, err := AssembleRefLeafTestSource(tc, compileOnly, pkgName, docTestRoot, rootImport, RefRootPkgName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", tc.Path, err)
	}
	testPath := filepath.Join(leafDir, TestFileName(tc))
	if err := WriteFormattedGo(testPath, src); err != nil {
		return "", err
	}
	return testPath, nil
}

// WriteFormattedGo formats src with imports.Process and writes atomically when changed.
func WriteFormattedGo(path, src string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	res, err := imports.Process(path, []byte(src), nil)
	if err != nil {
		_ = os.WriteFile(path, []byte(src), 0644)
		return fmt.Errorf("format imports for %s: %w", path, err)
	}
	existing, _ := os.ReadFile(path)
	if string(existing) == string(res) {
		return nil
	}
	tmpFile, err := os.CreateTemp(filepath.Dir(path), ".doctest-gen-*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	if _, err := tmpFile.Write(res); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		if !isCrossDeviceRename(err) {
			os.Remove(tmpPath)
			return err
		}
		if writeErr := os.WriteFile(path, res, 0644); writeErr != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("rename (cross-device): %w; write fallback: %v", err, writeErr)
		}
		os.Remove(tmpPath)
	}
	return nil
}

// CacheMappingGenRefRoot returns a cache gen root isolated from classic mapping-gen.
func CacheMappingGenRefRoot(absDoctestDir string) (string, string, error) {
	cacheDir, err := CacheHome()
	if err != nil {
		return "", "", err
	}
	absModRoot, _ := MappingGenRoot(absDoctestDir)
	mappingRoot := filepath.Join(cacheDir, "doctest", "mapping-gen-ref", absModRoot)
	return mappingRoot, absModRoot, nil
}
