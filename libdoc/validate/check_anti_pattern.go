package validate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func checkFileAntiPatterns(path string, content string) []error {
	if strings.Contains(path, "go-test-cache/env-getenv") {
		return nil
	}

	goCode := extractFinalGoBlock(content)
	if goCode == "" {
		return nil
	}

	var violations []error

	// Check for bare DOCTEST_ROOT / DOCTEST_CASE identifiers (not d.DOCTEST_*)
	if msg := checkBareDoctestRoots(path, goCode); msg != "" {
		violations = append(violations, fmt.Errorf("%s", msg))
	}

	fset := token.NewFileSet()
	file, parseErr := parser.ParseFile(fset, path+".go", "package testcase\n"+goCode, parser.ParseComments)
	if parseErr != nil {
		return violations
	}

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
					"%s: anti-pattern: do not combine a <contains> template with assert.Contains(); prefer assert.Output(t, actual, `` + `<contains>...</contains>`)",
					path,
				))
			}
			if msg := doctestSessionIDEnvReadMessage(node); msg != "" {
				violations = append(violations, fmt.Errorf("%s: %s", path, msg))
			}
			if msg := parallelUnsafeCallMessage(node); msg != "" {
				violations = append(violations, fmt.Errorf("%s: %s", path, msg))
			}
		case *ast.AssignStmt:
			if msg := parallelUnsafeStdioReassignMessage(node); msg != "" {
				violations = append(violations, fmt.Errorf("%s: %s", path, msg))
			}
		}
		return true
	})

	return violations
}

// parallelUnsafeCallMessage flags process-global env/cwd mutation APIs that
// race under t.Parallel() in suite harnesses.
func parallelUnsafeCallMessage(call *ast.CallExpr) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	recv, ok := sel.X.(*ast.Ident)
	if !ok {
		return ""
	}
	method := sel.Sel.Name
	switch recv.Name {
	case "os":
		switch method {
		case "Setenv", "Unsetenv", "Clearenv":
			return fmt.Sprintf(
				"anti-pattern: os.%s in harness under t.Parallel — set child cmd.Env only (or UseCLI subprocess); never mutate process env for leaf isolation",
				method,
			)
		case "Chdir":
			return "anti-pattern: os.Chdir in harness under t.Parallel — use absolute paths from d / t.TempDir() and child cmd.Dir; never change process cwd"
		}
	case "syscall":
		switch method {
		case "Setenv", "Unsetenv":
			return fmt.Sprintf(
				"anti-pattern: syscall.%s in harness under t.Parallel — set child cmd.Env only; never mutate process env for leaf isolation",
				method,
			)
		}
	case "t":
		switch method {
		case "Setenv":
			return "anti-pattern: t.Setenv in harness under t.Parallel — set child cmd.Env only; t.Setenv mutates process env and races parallel leaves"
		case "Chdir":
			return "anti-pattern: t.Chdir in harness under t.Parallel — use absolute paths and child cmd.Dir; t.Chdir mutates process cwd and races parallel leaves"
		}
	}
	return ""
}

// parallelUnsafeStdioReassignMessage flags assignment to os.Stdout / os.Stderr /
// os.Stdin (reads and Fprint to them are OK).
func parallelUnsafeStdioReassignMessage(assign *ast.AssignStmt) string {
	for _, lhs := range assign.Lhs {
		sel, ok := lhs.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pkg, ok := sel.X.(*ast.Ident)
		if !ok || pkg.Name != "os" {
			continue
		}
		switch sel.Sel.Name {
		case "Stdout", "Stderr", "Stdin":
			return fmt.Sprintf(
				"anti-pattern: os.%s in harness under t.Parallel — do not reassign process stdio; inject writers (opts/req) or fmt.Fprint(os.Stdout, …) without reassignment",
				sel.Sel.Name,
			)
		}
	}
	return ""
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

func doctestSessionIDEnvReadMessage(call *ast.CallExpr) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return ""
	}
	method := sel.Sel.Name
	switch pkg.Name {
	case "os":
		if method != "Getenv" && method != "LookupEnv" {
			return ""
		}
	case "syscall":
		if method != "Getenv" {
			return ""
		}
	default:
		return ""
	}
	if len(call.Args) == 0 {
		return ""
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	if stringLiteralValue(lit.Value) != "DOCTEST_SESSION_ID" {
		return ""
	}
	switch pkg.Name {
	case "os":
		return "anti-pattern: read DOCTEST_SESSION_ID via os." + method + " — use the doctest-injected variable DOCTEST_SESSION_ID directly (not an environment variable in harness code)"
	default:
		return "anti-pattern: read DOCTEST_SESSION_ID via syscall.Getenv — use the doctest-injected variable DOCTEST_SESSION_ID directly"
	}
}

// checkBareDoctestRoots scans source code for bare DOCTEST_ROOT, DOCTEST_CASE,
// or DOCTEST_SESSION_ID identifiers that aren't prefixed with d. (i.e. the old
// classic free-var style). In unified/ref gen mode, these are struct fields on
// d *session.Doctest and must be accessed as d.DOCTEST_ROOT, not bare.
func checkBareDoctestRoots(path, src string) string {
	// Check for bare identifiers by looking for the tokens NOT preceded by "d.".
	// Use a simple token-based scan on the raw source before AST parsing.
	for _, tok := range []string{"DOCTEST_ROOT", "DOCTEST_CASE", "DOCTEST_SESSION_ID"} {
		idx := 0
		for {
			pos := strings.Index(src[idx:], tok)
			if pos < 0 {
				break
			}
			absPos := idx + pos
			// Check if preceded by "d." — skip those (they're already correct).
			if absPos >= 2 && src[absPos-2:absPos] == "d." {
				idx = absPos + len(tok)
				continue
			}
			// Also skip if it's part of a string literal or comment.
			// Simple heuristic: if there's a // on this line before the token, it's a comment.
			lineStart := strings.LastIndexByte(src[:absPos], '\n')
			if lineStart < 0 {
				lineStart = 0
			}
			line := src[lineStart:absPos]
			if strings.Contains(line, "//") || strings.Contains(line, "`") {
				idx = absPos + len(tok)
				continue
			}
			return fmt.Sprintf("%s: %s used without d. prefix — add d *session.Doctest param and use d.%s instead of bare %s",
				path, tok, tok, tok)
		}
	}
	return ""
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
