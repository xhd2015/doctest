package core

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
)

const sessionImportPath = "github.com/xhd2015/doctest/session"

func AssembleTestSource(tc TreeCase, compileOnly bool, pkgName string, docTestRoot string) (string, error) {
	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteString("\n\n")

	imports := collectImports(tc.SetupFiles, tc.AssertFile.GoBlock)
	for _, pkg := range []string{"testing", "syscall", sessionImportPath} {
		if _, ok := imports[pkg]; !ok {
			imports[pkg] = &ImportSpec{Path: pkg}
		}
	}
	if len(imports) > 0 {
		buf.WriteString("import (\n")
		importList := make([]*ImportSpec, 0, len(imports))
		for _, spec := range imports {
			importList = append(importList, spec)
		}
		sort.Slice(importList, func(i, j int) bool { return importList[i].Path < importList[j].Path })
		for _, spec := range importList {
			if spec.Name != "" {
				buf.WriteString("\t" + spec.Name + " \"" + spec.Path + "\"\n")
			} else {
				buf.WriteString("\t\"" + spec.Path + "\"\n")
			}
		}
		buf.WriteString(")\n\n")
	}

	// Types, methods, consts/vars, and helpers at package level so methods can
	// implement interfaces and package helpers can reference shared vars
	// (e.g. uuidShape used by assertUUID). Go forbids methods on function-local types.
	// Session context is injected as d *session.Doctest (no package free DOCTEST_* vars).
	writePackageLevelTypesAndMethods(&buf, tc.SetupFiles, tc.AssertFile.GoBlock)
	writePackageLevelConstVars(&buf, tc.SetupFiles, tc.AssertFile.GoBlock)
	writePackageLevelHelpers(&buf, tc.SetupFiles, tc.AssertFile.GoBlock)

	buf.WriteString("func ")
	buf.WriteString(TestFuncName(tc))
	buf.WriteString("(t *testing.T) {\n")

	writeDoctestDConstruct(&buf, docTestRoot, tc.Path)

	var run *FuncSnippet
	if len(tc.SetupFiles) > 0 && tc.SetupFiles[0].GoBlock != nil && tc.SetupFiles[0].GoBlock.Run != nil {
		runCopy := *tc.SetupFiles[0].GoBlock.Run
		run = &runCopy
	}
	if run == nil {
		return "", fmt.Errorf("missing Run(t *testing.T, req *Request) (*Response, error) in setup chain")
	}
	run.Params = ensureDoctestParam(run.Params)
	writeFuncClosure(&buf, "run", *run)
	buf.WriteString("\tRun := run\n")

	buf.WriteString("\treq := &Request{}\n")
	writeSetupCalls(&buf, tc.SetupFiles)
	assertFn := *tc.AssertFile.GoBlock.Assert
	assertFn.Params = ensureDoctestParam(assertFn.Params)
	writeFuncClosure(&buf, "assert", assertFn)

	helperNames := collectHelperNames(tc.SetupFiles, tc.AssertFile.GoBlock)
	buf.WriteString("\t_ = Run\n")
	for _, name := range helperNames {
		buf.WriteString(fmt.Sprintf("\t_ = %s\n", name))
	}

	if compileOnly {
		buf.WriteString("\t// compileOnly\n")
		buf.WriteString("\t_ = d\n")
		buf.WriteString("\t_ = req\n")
		buf.WriteString("\t_ = run\n")
		buf.WriteString("\t_ = assert\n")
		buf.WriteString("\tvar resp *Response\n")
		buf.WriteString("\tvar runErr error\n")
		buf.WriteString("\t_ = resp\n")
		buf.WriteString("\t_ = runErr\n")
		buf.WriteString("}\n")
		return buf.String(), nil
	}
	buf.WriteString("\tresp, runErr := run(t, d, req)\n")
	buf.WriteString("\tassert(t, d, req, resp, runErr)\n")
	buf.WriteString("}\n")
	return buf.String(), nil
}

// writeDoctestDConstruct emits sid load + d := &session.Doctest{ROOT, CASE, SESSION_ID}.
// ROOT/CASE are baked absolute paths; SESSION_ID comes from the env at run time.
// d is declared in the enclosing test function scope for Setup/Run/Assert call sites.
func writeDoctestDConstruct(buf *strings.Builder, docTestRoot, casePath string) {
	escapedRoot := escapeRawString(docTestRoot)
	caseAbs := docTestRoot
	if casePath != "" {
		caseAbs = filepath.Join(docTestRoot, casePath)
	}
	escapedCase := escapeRawString(caseAbs)

	buf.WriteString("\tvar sid string\n")
	buf.WriteString("\t{\n")
	buf.WriteString("\t\ts, ok := syscall.Getenv(\"DOCTEST_SESSION_ID\")\n")
	buf.WriteString("\t\tif !ok || s == \"\" {\n")
	buf.WriteString("\t\t\tt.Fatalf(\"DOCTEST_SESSION_ID not set\")\n")
	buf.WriteString("\t\t}\n")
	buf.WriteString("\t\tsid = s\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\td := &session.Doctest{\n")
	buf.WriteString("\t\tDOCTEST_ROOT:       `" + escapedRoot + "`,\n")
	buf.WriteString("\t\tDOCTEST_CASE:       `" + escapedCase + "`,\n")
	buf.WriteString("\t\tDOCTEST_SESSION_ID: sid,\n")
	buf.WriteString("\t}\n")
}

func escapeRawString(s string) string {
	return strings.ReplaceAll(s, "`", "`+\"`\"+`")
}

// ensureDoctestParam inserts `_ *session.Doctest` after t when the author omitted it.
// If a param of type *session.Doctest is already present (any name), params are unchanged.
func ensureDoctestParam(params string) string {
	if strings.Contains(params, "*session.Doctest") {
		return params
	}
	trimmed := strings.TrimSpace(params)
	const tParam = "t *testing.T"
	if strings.HasPrefix(trimmed, tParam) {
		rest := strings.TrimSpace(trimmed[len(tParam):])
		rest = strings.TrimPrefix(rest, ",")
		rest = strings.TrimSpace(rest)
		if rest == "" {
			return tParam + ", _ *session.Doctest"
		}
		return tParam + ", _ *session.Doctest, " + rest
	}
	if trimmed == "" {
		return tParam + ", _ *session.Doctest"
	}
	return tParam + ", _ *session.Doctest, " + trimmed
}

func importKey(spec ImportSpec) string {
	if spec.Name != "" {
		return spec.Name + "\x00" + spec.Path
	}
	return spec.Path
}

func collectImports(setupFiles []SetupDocument, assertBlock GoBlock) map[string]*ImportSpec {
	imports := make(map[string]*ImportSpec)
	for _, doc := range setupFiles {
		if doc.GoBlock == nil {
			continue
		}
		for _, spec := range doc.GoBlock.Imports {
			if spec.Path == "" {
				continue
			}
			key := importKey(spec)
			if _, ok := imports[key]; !ok {
				imports[key] = &ImportSpec{Name: spec.Name, Path: spec.Path}
			}
		}
	}
	for _, spec := range assertBlock.Imports {
		if spec.Path == "" {
			continue
		}
		key := importKey(spec)
		if _, ok := imports[key]; !ok {
			imports[key] = &ImportSpec{Name: spec.Name, Path: spec.Path}
		}
	}
	return imports
}

func collectHelperNames(setupFiles []SetupDocument, assertBlock GoBlock) []string {
	seen := make(map[string]bool)
	var names []string
	for _, doc := range setupFiles {
		if doc.GoBlock == nil {
			continue
		}
		for _, h := range doc.GoBlock.Helpers {
			if !seen[h.Name] {
				seen[h.Name] = true
				names = append(names, h.Name)
			}
		}
	}
	for _, h := range assertBlock.Helpers {
		if !seen[h.Name] {
			seen[h.Name] = true
			names = append(names, h.Name)
		}
	}
	return names
}

func writePackageLevelTypesAndMethods(buf *strings.Builder, setupFiles []SetupDocument, assertBlock GoBlock) {
	// Deduplicate type decls by content (same harness types appear on every setup hop).
	seenTypes := make(map[string]bool)
	var typeDecls []string
	typeNames := make(map[string]bool)
	var methods []FuncSnippet
	seenMethod := make(map[string]bool)

	collect := func(block *GoBlock) {
		if block == nil {
			return
		}
		for name := range block.Types {
			typeNames[name] = true
		}
		for _, decl := range block.TypeDecls {
			if !seenTypes[decl] {
				seenTypes[decl] = true
				typeDecls = append(typeDecls, decl)
			}
		}
		for _, m := range block.Methods {
			key := m.Recv + " " + m.Name
			if !seenMethod[key] {
				seenMethod[key] = true
				methods = append(methods, m)
			}
		}
	}
	for _, doc := range setupFiles {
		collect(doc.GoBlock)
	}
	collect(&assertBlock)

	for _, decl := range sortTypeDecls(typeDecls, typeNames) {
		buf.WriteString(strings.TrimSpace(decl))
		buf.WriteString("\n")
	}
	if len(typeDecls) > 0 {
		buf.WriteString("\n")
	}
	for _, m := range methods {
		writeMethodDecl(buf, m)
	}
	if len(methods) > 0 {
		buf.WriteString("\n")
	}
}

func writeMethodDecl(buf *strings.Builder, fn FuncSnippet) {
	// Prefer ClosureResults: correctly parenthesizes multi-named results that
	// share one type (e.g. "(open []T, send []T)" not "open []T, send []T").
	results := fn.ClosureResults
	if strings.TrimSpace(results) == "" {
		results = fn.Results
	}
	if strings.TrimSpace(results) == "" {
		results = fn.ResultTypes
	}
	if strings.TrimSpace(results) != "" {
		results = " " + results
	}
	buf.WriteString(fmt.Sprintf("func (%s) %s(%s)%s %s\n", fn.Recv, fn.Name, fn.Params, results, fn.Body))
}

// writePackageLevelConstVars emits const/var decls at package scope so package
// helpers and methods can reference them (e.g. uuidShape, defaultAllowEmail).
func writePackageLevelConstVars(buf *strings.Builder, setupFiles []SetupDocument, assertBlock GoBlock) {
	seen := make(map[string]bool)
	emit := func(decls []string) {
		for _, decl := range decls {
			key := strings.TrimSpace(decl)
			if key == "" || seen[key] {
				continue
			}
			seen[key] = true
			buf.WriteString(key)
			buf.WriteString("\n")
		}
	}
	for _, doc := range setupFiles {
		if doc.GoBlock == nil {
			continue
		}
		emit(doc.GoBlock.Consts)
		emit(doc.GoBlock.VarDecls)
	}
	emit(assertBlock.Consts)
	emit(assertBlock.VarDecls)
	if len(seen) > 0 {
		buf.WriteString("\n")
	}
}

func writeConstVarDecls(buf *strings.Builder, setupFiles []SetupDocument) {
	for _, doc := range setupFiles {
		if doc.GoBlock == nil {
			continue
		}
		writeConstVarDeclsForBlock(buf, *doc.GoBlock)
	}
}

func writeConstVarDeclsForBlock(buf *strings.Builder, block GoBlock) {
	for _, decl := range block.Consts {
		writeIndented(buf, decl)
	}
	for _, decl := range block.VarDecls {
		writeIndented(buf, decl)
	}
}

// writeTypeConstVarDecls kept for any callers/tests that expect the old combined emit.
func writeTypeConstVarDecls(buf *strings.Builder, setupFiles []SetupDocument) {
	for _, doc := range setupFiles {
		if doc.GoBlock == nil {
			continue
		}
		writeTypeConstVarDeclsForBlock(buf, *doc.GoBlock)
	}
}

func writeTypeConstVarDeclsForBlock(buf *strings.Builder, block GoBlock) {
	for _, decl := range sortTypeDecls(block.TypeDecls, block.Types) {
		writeIndented(buf, decl)
	}
	for _, decl := range block.Consts {
		writeIndented(buf, decl)
	}
	for _, decl := range block.VarDecls {
		writeIndented(buf, decl)
	}
}

func writeAllHelpers(buf *strings.Builder, setupFiles []SetupDocument, assertBlock GoBlock) {
	// Legacy: emit helpers as closures inside the test function.
	var helpers []FuncSnippet
	for _, doc := range setupFiles {
		if doc.GoBlock == nil {
			continue
		}
		helpers = append(helpers, doc.GoBlock.Helpers...)
	}
	helpers = append(helpers, assertBlock.Helpers...)
	for _, helper := range sortHelpers(helpers) {
		writeFuncClosure(buf, helper.Name, helper)
	}
}

// writePackageLevelHelpers emits non-method helpers as real package functions
// so package-level methods can call them (e.g. containsArg from fakeRunner.Exec).
func writePackageLevelHelpers(buf *strings.Builder, setupFiles []SetupDocument, assertBlock GoBlock) {
	seen := make(map[string]bool)
	var helpers []FuncSnippet
	for _, doc := range setupFiles {
		if doc.GoBlock == nil {
			continue
		}
		for _, h := range doc.GoBlock.Helpers {
			if !seen[h.Name] {
				seen[h.Name] = true
				helpers = append(helpers, h)
			}
		}
	}
	for _, h := range assertBlock.Helpers {
		if !seen[h.Name] {
			seen[h.Name] = true
			helpers = append(helpers, h)
		}
	}
	for _, helper := range sortHelpers(helpers) {
		writePackageFunc(buf, helper)
	}
	if len(helpers) > 0 {
		buf.WriteString("\n")
	}
}

func writePackageFunc(buf *strings.Builder, fn FuncSnippet) {
	// Prefer ClosureResults so multi-named shared-type results are parenthesized.
	results := fn.ClosureResults
	if strings.TrimSpace(results) == "" {
		results = fn.Results
	}
	if strings.TrimSpace(results) == "" {
		results = fn.ResultTypes
	}
	if strings.TrimSpace(results) != "" {
		results = " " + results
	}
	buf.WriteString(fmt.Sprintf("func %s(%s)%s %s\n", fn.Name, fn.Params, results, fn.Body))
}

// sortHelpers topologically orders helpers so that a closure referencing
// another helper is emitted after its callee. Top-level funcs allow forward
// references; func literals (closures) do not, so without reordering a helper
// defined before one it calls fails to compile. On a cycle (mutual recursion)
// it falls back to source order.
func sortHelpers(helpers []FuncSnippet) []FuncSnippet {
	type depInfo struct {
		helper FuncSnippet
		deps   []int // indices of helpers this one calls
	}
	nameIdx := make(map[string]int)
	infos := make([]depInfo, len(helpers))
	for i, h := range helpers {
		infos[i].helper = h
		nameIdx[h.Name] = i
	}
	for i, h := range helpers {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "", "package p; func _() "+h.Body, 0)
		if err != nil {
			continue
		}
		seen := make(map[int]bool)
		ast.Inspect(file, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if j, ok := nameIdx[id.Name]; ok && j != i && !seen[j] {
				infos[i].deps = append(infos[i].deps, j)
				seen[j] = true
			}
			return true
		})
	}
	inDegree := make(map[int]int)
	adj := make(map[int][]int)
	for i, info := range infos {
		for _, j := range info.deps {
			adj[j] = append(adj[j], i)
			inDegree[i]++
		}
	}
	var sorted []FuncSnippet
	var queue []int
	for i := range infos {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		sort.Ints(queue) // prefer source order among ready nodes
		i := queue[0]
		queue = queue[1:]
		sorted = append(sorted, infos[i].helper)
		for _, j := range adj[i] {
			inDegree[j]--
			if inDegree[j] == 0 {
				queue = append(queue, j)
			}
		}
	}
	if len(sorted) == len(helpers) {
		return sorted
	}
	return helpers
}

func sortTypeDecls(decls []string, types map[string]bool) []string {
	type depInfo struct {
		decl string
		deps []string
	}
	var infos []depInfo
	declNames := make(map[string]int)
	for _, decl := range decls {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "", "package p\n"+decl, 0)
		if err != nil {
			continue
		}
		var typeName string
		for _, d := range file.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				typeName = ts.Name.Name
			}
		}
		var deps []string
		ast.Inspect(file, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if id.Name == typeName {
				return true
			}
			if types[id.Name] {
				deps = append(deps, id.Name)
			}
			return true
		})
		infos = append(infos, depInfo{decl: decl, deps: deps})
		declNames[typeName] = len(infos) - 1
	}

	inDegree := make(map[int]int)
	adj := make(map[int][]int)
	for i, info := range infos {
		for _, dep := range info.deps {
			if j, ok := declNames[dep]; ok {
				adj[j] = append(adj[j], i)
				inDegree[i]++
			}
		}
	}

	var sorted []string
	var queue []int
	for i := range infos {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		i := queue[0]
		queue = queue[1:]
		sorted = append(sorted, infos[i].decl)
		for _, j := range adj[i] {
			inDegree[j]--
			if inDegree[j] == 0 {
				queue = append(queue, j)
			}
		}
	}
	if len(sorted) == len(decls) {
		return sorted
	}
	return decls
}

func writeSetupCalls(buf *strings.Builder, setupFiles []SetupDocument) {
	for i, doc := range setupFiles {
		if doc.GoBlock == nil || doc.GoBlock.Setup == nil {
			continue
		}
		name := fmt.Sprintf("setup%d", i)
		fn := *doc.GoBlock.Setup
		fn.Params = ensureDoctestParam(fn.Params)
		writeFuncClosure(buf, name, fn)
		buf.WriteString(fmt.Sprintf("\tif err := %s(t, d, req); err != nil {\n", name))
		buf.WriteString(fmt.Sprintf("\t\tt.Fatalf(\"%s failed: %%v\", err)\n", escapeString(doc.Path)))
		buf.WriteString("\t}\n")
	}
}

func writeFuncClosure(buf *strings.Builder, name string, fn FuncSnippet) {
	// Prefer the closure-rendered results (names preserved, correctly
	// parenthesized). Fall back to ResultTypes for manually-built snippets
	// that only set type-only results, then to Results.
	results := fn.ClosureResults
	if strings.TrimSpace(results) == "" {
		results = fn.ResultTypes
	}
	if strings.TrimSpace(results) == "" {
		results = fn.Results
	}
	if strings.TrimSpace(results) != "" {
		results = " " + results
	}
	buf.WriteString(fmt.Sprintf("\t%s := func(%s)%s %s\n", name, fn.Params, results, fn.Body))
}

func writeIndented(buf *strings.Builder, s string) {
	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		buf.WriteString("\t")
		buf.WriteString(line)
		buf.WriteString("\n")
	}
}

func escapeString(s string) string {
	return strings.ReplaceAll(s, "\\", "\\\\")
}
