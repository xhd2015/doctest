package core

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// formatGeneratedGo reconciles imports on assembled source without go env /
// imports.Process, then prints the AST. It does not call go/format.Source and
// never auto-adds stdlib (or any package) from selector usage — authors must
// import packages their own code names. Assemble may over-import ancestors;
// unused imports are pruned. On parse failure the original bytes are returned
// unchanged so a compiling (if oddly formatted) package can still be written.
func formatGeneratedGo(_ string, src []byte) ([]byte, error) {
	fixed, err := reconcileGeneratedImports(string(src))
	if err != nil {
		// No format.Source fallback — leave source as assembled.
		return src, nil
	}
	return []byte(fixed), nil
}

func reconcileGeneratedImports(src string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "generated.go", src, parser.ParseComments)
	if err != nil {
		return "", err
	}

	used := map[string]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if id, ok := sel.X.(*ast.Ident); ok && id != nil {
			used[id.Name] = true
		}
		return true
	})

	// Drop unused imports (keep blank imports).
	// DeleteNamedImport needs the AST name ("" if unaliased, not path base).
	type impRef struct {
		path      string
		localName string // name used in code (alias or path base)
		astName   string // astutil name: "" unaliased, else alias
		blank     bool
	}
	var refs []impRef
	for _, imp := range file.Imports {
		p, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}
		blank := imp.Name != nil && imp.Name.Name == "_"
		astName := ""
		if imp.Name != nil {
			astName = imp.Name.Name
		}
		refs = append(refs, impRef{
			path:      p,
			localName: importLocalName(imp, p),
			astName:   astName,
			blank:     blank,
		})
	}
	for _, r := range refs {
		if r.blank {
			continue
		}
		if r.astName == "." {
			continue
		}
		if used[r.localName] {
			continue
		}
		astutil.DeleteNamedImport(fset, file, r.astName, r.path)
	}

	var buf bytes.Buffer
	// format.Node prints the AST; this is not go/format.Source and is only
	// used to emit the reconciled file after import prune.
	if err := format.Node(&buf, fset, file); err != nil {
		return "", fmt.Errorf("format after import reconcile: %w", err)
	}
	return buf.String(), nil
}

func importLocalName(imp *ast.ImportSpec, importPath string) string {
	if imp.Name != nil {
		return imp.Name.Name
	}
	// path/filepath → filepath; github.com/foo/bar → bar
	base := path.Base(importPath)
	if base == "/" || base == "." || base == "" {
		return importPath
	}
	// strip version suffixes rarely present in imports
	if i := strings.Index(base, "."); i > 0 && !strings.Contains(importPath, "/") {
		// stdlib single segment
		return base
	}
	return base
}
