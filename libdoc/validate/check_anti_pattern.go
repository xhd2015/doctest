package validate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func checkFileAntiPatterns(path string, content string) []error {
	goCode := extractFinalGoBlock(content)
	if goCode == "" {
		return nil
	}

	fset := token.NewFileSet()
	file, parseErr := parser.ParseFile(fset, path+".go", "package testcase\n"+goCode, parser.ParseComments)
	if parseErr != nil {
		return nil
	}

	var violations []error
	containsTemplateVars := collectContainsTemplateVars(file)
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.BasicLit:
			if node.Kind == token.STRING && containsEmbeddedGoProgram(node.Value) {
				violations = append(violations, fmt.Errorf(
					"%s: anti-pattern: raw Go code embedded in string literal (contains 'package main' and 'func main()') — import the package and call its functions directly",
					path,
				))
			}
		case *ast.CallExpr:
			if isGoTestShellOut(node) {
				violations = append(violations, fmt.Errorf(
					"%s: anti-pattern: shelling out to 'go test' — call the function under test directly from Run",
					path,
				))
			}
			if isDoubleContainsOutputAssert(node, containsTemplateVars) {
				violations = append(violations, fmt.Errorf(
					"%s: anti-pattern: do not combine a <contains> template with assert.Contains(); prefer assert.Output(t, actual, `<contains>...</contains>`)",
					path,
				))
			}
		}
		return true
	})

	return violations
}

func collectContainsTemplateVars(file *ast.File) map[string]bool {
	vars := make(map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for i, rhs := range assign.Rhs {
			if !isAssertParseContainsTemplate(rhs) {
				continue
			}
			if i >= len(assign.Lhs) {
				continue
			}
			if ident, ok := assign.Lhs[i].(*ast.Ident); ok {
				vars[ident.Name] = true
			}
		}
		return true
	})
	return vars
}

func extractFinalGoBlock(content string) string {
	blocks := findGoBlocks(content)
	if len(blocks) == 0 {
		return ""
	}
	return blocks[len(blocks)-1].code
}

type mdBlock struct {
	code string
	end  int
}

func findGoBlocks(content string) []mdBlock {
	var blocks []mdBlock
	i := 0
	for {
		start := strings.Index(content[i:], "```go")
		if start < 0 {
			return blocks
		}
		start += i
		lineEnd := strings.IndexByte(content[start:], '\n')
		if lineEnd < 0 {
			return blocks
		}
		codeStart := start + lineEnd + 1
		close := strings.Index(content[codeStart:], "```")
		if close < 0 {
			return blocks
		}
		close += codeStart
		code := content[codeStart:close]
		end := close + len("```")
		if end < len(content) && content[end] == '\r' {
			end++
		}
		if end < len(content) && content[end] == '\n' {
			end++
		}
		blocks = append(blocks, mdBlock{code: code, end: end})
		i = end
	}
}

func containsEmbeddedGoProgram(lit string) bool {
	if strings.HasPrefix(lit, "`") && strings.HasSuffix(lit, "`") {
		s := lit[1 : len(lit)-1]
		return strings.Contains(s, "package main") && strings.Contains(s, "func main()")
	}
	if strings.HasPrefix(lit, `"`) && strings.HasSuffix(lit, `"`) {
		s := lit[1 : len(lit)-1]
		return strings.Contains(s, "package main") && strings.Contains(s, "func main()")
	}
	return false
}

func isDoubleContainsOutputAssert(call *ast.CallExpr, containsTemplateVars map[string]bool) bool {
	if !isSelectorNamed(call.Fun, "Match") {
		return false
	}
	if !callHasContainsOption(call) {
		return false
	}
	if len(call.Args) == 0 {
		return false
	}
	if isAssertParseContainsTemplate(call.Args[0]) {
		return true
	}
	if ident, ok := call.Args[0].(*ast.Ident); ok {
		return containsTemplateVars[ident.Name]
	}
	return false
}

func callHasContainsOption(call *ast.CallExpr) bool {
	for _, arg := range call.Args[1:] {
		if inner, ok := arg.(*ast.CallExpr); ok && isSelectorNamed(inner.Fun, "Contains") {
			return true
		}
	}
	return false
}

func isAssertParseContainsTemplate(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	if !isSelectorNamed(call.Fun, "Parse") && !isSelectorNamed(call.Fun, "MustParse") {
		return false
	}
	if len(call.Args) == 0 {
		return false
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	text := stringLiteralValue(lit.Value)
	return strings.Contains(text, "<contains>")
}

func isSelectorNamed(expr ast.Expr, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == name
}

func stringLiteralValue(lit string) string {
	if len(lit) < 2 {
		return lit
	}
	if strings.HasPrefix(lit, "`") && strings.HasSuffix(lit, "`") {
		return lit[1 : len(lit)-1]
	}
	if strings.HasPrefix(lit, `"`) && strings.HasSuffix(lit, `"`) {
		return strings.Trim(lit, `"`)
	}
	return lit
}

func isGoTestShellOut(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	if sel.Sel.Name != "Command" || ident.Name != "exec" {
		return false
	}
	if len(call.Args) < 2 {
		return false
	}
	firstArg, ok1 := call.Args[0].(*ast.BasicLit)
	secondArg, ok2 := call.Args[1].(*ast.BasicLit)
	if !ok1 || !ok2 {
		return false
	}
	if firstArg.Kind != token.STRING || secondArg.Kind != token.STRING {
		return false
	}
	firstVal := strings.Trim(firstArg.Value, `"`)
	secondVal := strings.Trim(secondArg.Value, `"`)
	return firstVal == "go" && secondVal == "test"
}
