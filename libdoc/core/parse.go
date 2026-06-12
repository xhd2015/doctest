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
	if strings.TrimSpace(content[last.end:]) != "" {
		return GoBlock{}, fmt.Errorf("%s: go block must be final content", path)
	}
	return GoBlock{SourcePath: path, Code: last.code}, nil
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
		if v := rules.CheckSetupBodyNotStub(block.Setup.Body, path); v != nil {
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
	return AssertDocument{Path: path, GoBlock: block}, nil
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
						block.Imports = append(block.Imports, strings.Trim(is.Path.Value, "\""))
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
	if d.Type.Results != nil {
		results = resultsString(fset, d.Type.Results)
	}
	return FuncSnippet{
		Name:    d.Name.Name,
		Params:  params,
		Results: results,
		Body:    nodeString(fset, d.Body),
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
