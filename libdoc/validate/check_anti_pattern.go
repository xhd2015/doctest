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
		}
		return true
	})

	return violations
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
