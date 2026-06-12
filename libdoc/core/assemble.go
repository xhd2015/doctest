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
	imports["testing"] = true
	imports["os"] = true
	imports["path/filepath"] = true
	if len(imports) > 0 {
		buf.WriteString("import (\n")
		importList := make([]string, 0, len(imports))
		for pkg := range imports {
			importList = append(importList, pkg)
		}
		sort.Strings(importList)
		for _, pkg := range importList {
			buf.WriteString("\t\"" + pkg + "\"\n")
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

	writeHelperDecls(&buf, tc.SetupFiles)
	writeHelperDeclsForBlock(&buf, tc.AssertFile.GoBlock)

	var run *FuncSnippet
	for _, doc := range tc.SetupFiles {
		if doc.GoBlock != nil && doc.GoBlock.Run != nil {
			runCopy := *doc.GoBlock.Run
			run = &runCopy
		}
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

func collectImports(setupFiles []SetupDocument, assertBlock GoBlock) map[string]bool {
	imports := make(map[string]bool)
	for _, doc := range setupFiles {
		if doc.GoBlock == nil {
			continue
		}
		for _, pkg := range doc.GoBlock.Imports {
			if pkg != "" {
				imports[pkg] = true
			}
		}
	}
	for _, pkg := range assertBlock.Imports {
		if pkg != "" {
			imports[pkg] = true
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

func writeHelperDecls(buf *strings.Builder, setupFiles []SetupDocument) {
	for _, doc := range setupFiles {
		if doc.GoBlock == nil {
			continue
		}
		for _, helper := range doc.GoBlock.Helpers {
			writeFuncClosure(buf, helper.Name, helper)
		}
	}
}

func writeHelperDeclsForBlock(buf *strings.Builder, block GoBlock) {
	for _, helper := range block.Helpers {
		writeFuncClosure(buf, helper.Name, helper)
	}
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
	results := ""
	if strings.TrimSpace(fn.Results) != "" {
		results = " " + fn.Results
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
