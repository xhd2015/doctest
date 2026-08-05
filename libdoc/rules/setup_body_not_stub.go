package rules

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// vacuousSetupMsg tells authors to drop empty Setup rather than "implement" it.
// Organization-only nodes should use prose SETUP (no Go block) or omit Setup.
const vacuousSetupMsg = "vacuous func Setup (only return nil / blank assigns) — remove the Go code block (or omit Setup); organization-only nodes need no Setup"

// CheckSetupBodyNotStub reports vacuous non-root Setup bodies.
// Vacuous if the body is only:
//   - return nil, or
//   - zero or more blank-identifier assignments (_ = expr) then return nil
// Comments-only with return nil is still vacuous (comments are not statements).
func CheckSetupBodyNotStub(body string, path string) *Violation {
	return CheckVacuousSetup(body, path)
}

// CheckVacuousSetup is the preferred name for vacuous Setup body detection.
func CheckVacuousSetup(body string, path string) *Violation {
	body = strings.TrimSpace(body)
	// Body from go/printer is typically "{ ... }"; normalize to inner stmts.
	inner := body
	inner = strings.TrimPrefix(inner, "{")
	inner = strings.TrimSuffix(inner, "}")
	inner = strings.TrimSpace(inner)
	if inner == "" || inner == "return nil" {
		return &Violation{Path: path, Msg: vacuousSetupMsg}
	}

	// Re-wrap braces so multi-stmt bodies (_ = x; return nil) parse correctly.
	fset := token.NewFileSet()
	src := "package p\nfunc _() {\n" + inner + "\n}"
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil || len(file.Decls) == 0 {
		return nil
	}
	fd, ok := file.Decls[0].(*ast.FuncDecl)
	if !ok || fd.Body == nil {
		return nil
	}
	if isVacuousSetupBody(fd.Body.List) {
		return &Violation{Path: path, Msg: vacuousSetupMsg}
	}
	return nil
}

// isVacuousSetupBody reports whether stmts are only blank assigns + final return nil.
func isVacuousSetupBody(stmts []ast.Stmt) bool {
	if len(stmts) == 0 {
		return true
	}
	for i, stmt := range stmts {
		last := i == len(stmts)-1
		if last {
			return isReturnNil(stmt)
		}
		if !isBlankAssign(stmt) {
			return false
		}
	}
	return false
}

func isReturnNil(stmt ast.Stmt) bool {
	rs, ok := stmt.(*ast.ReturnStmt)
	if !ok || len(rs.Results) != 1 {
		return false
	}
	id, ok := rs.Results[0].(*ast.Ident)
	return ok && id.Name == "nil"
}

// isBlankAssign reports `_ = expr` or `_ , _ = ...` style blank-only assigns.
// Only pure blank-identifier LHS counts (discarding params without real work).
func isBlankAssign(stmt ast.Stmt) bool {
	as, ok := stmt.(*ast.AssignStmt)
	if !ok {
		return false
	}
	if len(as.Lhs) == 0 {
		return false
	}
	for _, lhs := range as.Lhs {
		id, ok := lhs.(*ast.Ident)
		if !ok || id.Name != "_" {
			return false
		}
	}
	return true
}
