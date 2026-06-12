package rules

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func CheckSetupBodyNotStub(body string, path string) *Violation {
	body = strings.TrimSpace(body)
	body = strings.TrimPrefix(body, "{")
	body = strings.TrimSuffix(body, "}")
	body = strings.TrimSpace(body)
	if body == "" || body == "return nil" {
		return &Violation{Path: path, Msg: "func Setup body must not be a stub (return nil) — implement the behavior described in this document"}
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", "package p\nfunc _() "+body, parser.ParseComments)
	if err != nil || len(file.Decls) == 0 {
		return nil
	}
	fd, ok := file.Decls[0].(*ast.FuncDecl)
	if !ok || fd.Body == nil {
		return nil
	}
	stmts := fd.Body.List
	if len(stmts) == 1 {
		if rs, ok := stmts[0].(*ast.ReturnStmt); ok {
			if len(rs.Results) == 1 {
				if id, ok := rs.Results[0].(*ast.Ident); ok && id.Name == "nil" {
					return &Violation{Path: path, Msg: "func Setup body must not be a stub (return nil) — implement the behavior described in this document"}
				}
			}
		}
	}
	return nil
}
