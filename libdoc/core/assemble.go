package core

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
)

func AssembleTestSource(tc TreeCase, compileOnly bool, pkgName string, docTestRoot string) (string, error) {
	var buf strings.Builder
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteString("\n\n")

	imports := collectImports(tc.SetupFiles, tc.AssertFile.GoBlock)
	for _, pkg := range []string{"testing", "os", "path/filepath", "syscall"} {
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

	buf.WriteString("func ")
	buf.WriteString(TestFuncName(tc))
	buf.WriteString("(t *testing.T) {\n")

	escapedRoot := strings.ReplaceAll(docTestRoot, "`", "`+\"`\"+`")
	buf.WriteString("\tconst DOCTEST_ROOT = `")
	buf.WriteString(escapedRoot)
	buf.WriteString("`\n")
	buf.WriteString("\tDOCTEST_SESSION_ID, __sessionOk := syscall.Getenv(\"DOCTEST_SESSION_ID\")\n")
	buf.WriteString("\tif !__sessionOk || DOCTEST_SESSION_ID == \"\" {\n")
	buf.WriteString("\t\tt.Fatalf(\"DOCTEST_SESSION_ID not set\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\t__origWd, __wdErr := os.Getwd()\n")
	buf.WriteString("\tif __wdErr != nil {\n")
	buf.WriteString("\t\tt.Fatal(__wdErr)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tdefer os.Chdir(__origWd)\n")
	if tc.Path != "" {
		buf.WriteString(fmt.Sprintf("\tif err := os.Chdir(filepath.Join(DOCTEST_ROOT, %q)); err != nil {\n", tc.Path))
	} else {
		buf.WriteString("\tif err := os.Chdir(DOCTEST_ROOT); err != nil {\n")
	}
	buf.WriteString("\t\tt.Fatal(err)\n")
	buf.WriteString("\t}\n\n")

	writeTypeConstVarDecls(&buf, tc.SetupFiles)
	writeTypeConstVarDeclsForBlock(&buf, tc.AssertFile.GoBlock)

	writeAllHelpers(&buf, tc.SetupFiles, tc.AssertFile.GoBlock)

	var run *FuncSnippet
	if len(tc.SetupFiles) > 0 && tc.SetupFiles[0].GoBlock != nil && tc.SetupFiles[0].GoBlock.Run != nil {
		runCopy := *tc.SetupFiles[0].GoBlock.Run
		run = &runCopy
	}
	if run == nil {
		return "", fmt.Errorf("missing Run(t *testing.T, req *Request) (*Response, error) in setup chain")
	}
	writeFuncClosure(&buf, "run", *run)
	buf.WriteString("\tRun := run\n")

	buf.WriteString("\treq := &Request{}\n")
	writeSetupCalls(&buf, tc.SetupFiles)
	writeFuncClosure(&buf, "assert", *tc.AssertFile.GoBlock.Assert)

	helperNames := collectHelperNames(tc.SetupFiles, tc.AssertFile.GoBlock)
	buf.WriteString("\t_ = Run\n")
	for _, name := range helperNames {
		buf.WriteString(fmt.Sprintf("\t_ = %s\n", name))
	}

	if compileOnly {
		buf.WriteString("\t// compileOnly\n")
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
	buf.WriteString("\tresp, runErr := run(t, req)\n")
	buf.WriteString("\tassert(t, req, resp, runErr)\n")
	buf.WriteString("}\n")
	return buf.String(), nil
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
		writeFuncClosure(buf, name, *doc.GoBlock.Setup)
		buf.WriteString(fmt.Sprintf("\tif err := %s(t, req); err != nil {\n", name))
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
