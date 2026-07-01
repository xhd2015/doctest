package core

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"

	"github.com/xhd2015/doctest/libdoc/rules"
)

func ExtractFinalGoBlock(path string, content string) (GoBlock, error) {
	blocks := findGoBlocks(content)
	if len(blocks) == 0 {
		return GoBlock{}, fmt.Errorf("%s: missing go block", path)
	}
	last := blocks[len(blocks)-1]
	if trailing := strings.TrimSpace(content[last.end:]); trailing != "" && !isAllowedTrailingContent(trailing) {
		return GoBlock{}, fmt.Errorf("%s: go block must be final content", path)
	}
	return GoBlock{SourcePath: path, Code: last.code}, nil
}

func isAllowedTrailingContent(s string) bool {
	rest := strings.TrimSpace(s)
	for rest != "" {
		if !strings.HasPrefix(rest, "<!--") {
			return false
		}
		end := strings.Index(rest, "-->")
		if end < 0 {
			return false
		}
		rest = strings.TrimSpace(rest[end+len("-->"):])
	}
	return true
}

func ParseDOCTESTDocument(path string, content string) (SetupDocument, error) {
	block, err := ExtractFinalGoBlock(path, content)
	if err != nil {
		return SetupDocument{}, err
	}
	if err := parseGoBlock(&block); err != nil {
		return SetupDocument{}, err
	}
	if block.Run != nil {
		if v := rules.CheckRunSignature(block.Run.Params, block.Run.Results, path); v != nil {
			return SetupDocument{}, fmt.Errorf("%s: %s", v.Path, v.Msg)
		}
	}
	return SetupDocument{Path: path, GoBlock: &block}, nil
}

func ParseSetupDocument(path string, content string) (SetupDocument, error) {
	blocks := findGoBlocks(content)
	if len(blocks) == 0 {
		return SetupDocument{Path: path}, nil
	}
	if len(blocks) > 1 {
		return SetupDocument{}, fmt.Errorf("%s: multiple go blocks are not allowed", path)
	}
	block, err := ExtractFinalGoBlock(path, content)
	if err != nil {
		return SetupDocument{}, err
	}
	if err := parseGoBlock(&block); err != nil {
		return SetupDocument{}, err
	}
	if block.Setup != nil {
		if v := rules.CheckSetupSignature(block.Setup.Params, block.Setup.Results, path); v != nil {
			return SetupDocument{}, fmt.Errorf("%s: %s", v.Path, v.Msg)
		}
	}
	if block.Run != nil {
		if v := rules.CheckRunSignature(block.Run.Params, block.Run.Results, path); v != nil {
			return SetupDocument{}, fmt.Errorf("%s: %s", v.Path, v.Msg)
		}
	}
	return SetupDocument{Path: path, GoBlock: &block}, nil
}

func ParseAssertDocument(path string, content string) (AssertDocument, error) {
	fm, body, err := ParseAssertFrontmatter(path, content)
	if err != nil {
		return AssertDocument{}, err
	}
	content = body
	blocks := findGoBlocks(content)
	if len(blocks) == 0 {
		return AssertDocument{}, fmt.Errorf("%s: missing go block", path)
	}
	if len(blocks) > 1 {
		return AssertDocument{}, fmt.Errorf("%s: multiple go blocks are not allowed", path)
	}
	block, err := ExtractFinalGoBlock(path, content)
	if err != nil {
		return AssertDocument{}, err
	}
	if err := parseGoBlock(&block); err != nil {
		return AssertDocument{}, err
	}
	if v := rules.CheckAssertExists(block.Assert != nil, path); v != nil {
		return AssertDocument{}, fmt.Errorf("%s: %s", v.Path, v.Msg)
	}
	if v := rules.CheckAssertSignature(block.Assert.Params, block.Assert.Results, path); v != nil {
		return AssertDocument{}, fmt.Errorf("%s: %s", v.Path, v.Msg)
	}
	return AssertDocument{Path: path, GoBlock: block, Frontmatter: fm}, nil
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

func parseGoBlock(block *GoBlock) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, block.SourcePath+".go", "package testcase\n"+block.Code, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("%s: invalid go: %w", block.SourcePath, err)
	}
	block.Types = make(map[string]bool)
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
		case token.IMPORT:
			for _, spec := range d.Specs {
				if is, ok := spec.(*ast.ImportSpec); ok {
					var name string
					if is.Name != nil {
						name = is.Name.Name
					}
					block.Imports = append(block.Imports, ImportSpec{Name: name, Path: strings.Trim(is.Path.Value, "\"")})
				}
			}
			case token.TYPE:
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						block.Types[ts.Name.Name] = true
					}
				}
				block.TypeDecls = append(block.TypeDecls, nodeString(fset, d))
			case token.VAR:
				block.VarDecls = append(block.VarDecls, nodeString(fset, d))
			case token.CONST:
				block.Consts = append(block.Consts, nodeString(fset, d))
			}
		case *ast.FuncDecl:
			fn := funcSnippet(fset, d)
			switch d.Name.Name {
			case "Setup":
				block.Setup = &fn
			case "Run":
				block.Run = &fn
			case "Assert":
				block.Assert = &fn
			default:
				block.Helpers = append(block.Helpers, fn)
			}
		}
	}
	return nil
}

func nodeString(fset *token.FileSet, n any) string {
	var buf bytes.Buffer
	_ = printer.Fprint(&buf, fset, n)
	return buf.String()
}

func funcSnippet(fset *token.FileSet, d *ast.FuncDecl) FuncSnippet {
	params := ""
	if d.Type.Params != nil {
		params = fieldsString(fset, d.Type.Params)
	}
	results := ""
	resultTypes := ""
	closureResults := ""
	if d.Type.Results != nil {
		results = resultsString(fset, d.Type.Results)
		resultTypes = resultTypesString(fset, d.Type.Results)
		closureResults = closureResultsString(fset, d.Type.Results)
	}
	return FuncSnippet{
		Name:           d.Name.Name,
		Params:         params,
		Results:        results,
		ResultTypes:    resultTypes,
		ClosureResults: closureResults,
		Body:           nodeString(fset, d.Body),
	}
}

func fieldsString(fset *token.FileSet, fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}
	var parts []string
	for _, field := range fields.List {
		typ := nodeString(fset, field.Type)
		if len(field.Names) == 0 {
			parts = append(parts, typ)
			continue
		}
		for _, name := range field.Names {
			parts = append(parts, name.Name+" "+typ)
		}
	}
	return strings.Join(parts, ", ")
}

func resultsString(fset *token.FileSet, fields *ast.FieldList) string {
	if fields == nil || len(fields.List) == 0 {
		return ""
	}
	s := fieldsString(fset, fields)
	if len(fields.List) > 1 {
		return "(" + s + ")"
	}
	return s
}

// resultTypesString renders result parameters with type names only (no identifiers).
// Multiple named results sharing one type, e.g. (port, alt int), must be parenthesized
// in func literals: (int, int) not port int, alt int.
func resultTypesString(fset *token.FileSet, fields *ast.FieldList) string {
	if fields == nil || len(fields.List) == 0 {
		return ""
	}
	var parts []string
	for _, field := range fields.List {
		typ := nodeString(fset, field.Type)
		n := len(field.Names)
		if n == 0 {
			n = 1
		}
		for i := 0; i < n; i++ {
			parts = append(parts, typ)
		}
	}
	s := strings.Join(parts, ", ")
	if len(parts) > 1 {
		return "(" + s + ")"
	}
	return s
}

// closureResultsString renders result parameters for a func literal: it keeps
// parameter names (so bodies that assign to named returns compile) and
// parenthesizes whenever there is more than one result or any named result,
// because bare multi-name results like "port, alt int" are invalid syntax
// outside parentheses. A single unnamed result (e.g. "int") stays bare.
func closureResultsString(fset *token.FileSet, fields *ast.FieldList) string {
	if fields == nil || len(fields.List) == 0 {
		return ""
	}
	s := fieldsString(fset, fields)
	total := 0
	anyNamed := false
	for _, f := range fields.List {
		if len(f.Names) == 0 {
			total++
		} else {
			total += len(f.Names)
			anyNamed = true
		}
	}
	if total > 1 || anyNamed {
		return "(" + s + ")"
	}
	return s
}
