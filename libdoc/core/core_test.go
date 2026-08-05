package core

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

)

func TestRootSetupRequiresRequestAndResponseTypes(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return nil, nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))

	_, err := DiscoverTreeCases(root)
	if err == nil {
		t.Fatal("expected missing Request/Response error")
	}
	if !strings.Contains(err.Error(), "Request") {
		t.Fatalf("expected Request error, got %v", err)
	}
}

func TestValidationReportsAllErrorsAtOnce(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return nil, nil }`))
	writeTreeFile(t, root, "SETUP.md", "# Setup\n\n```go\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }\n```\n")
	writeTreeFile(t, root, "leaf/SETUP.md", "# Setup\n\nProse only, no Go block.\n")
	writeTreeFile(t, root, "leaf/ASSERT.md", "# Assert\n\n```go\nfunc Check(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}\n```\n")

	_, err := DiscoverTreeCases(root)
	if err == nil {
		t.Fatal("expected validation error")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "must define type Request") && !strings.Contains(errStr, "type Response") {
		t.Fatalf("expected missing types error, got %q", errStr)
	}
	// Prose-only SETUP.md (no Go block) is allowed for organization-only nodes.
	if strings.Contains(errStr, "must have a Go code block") {
		t.Fatalf("prose-only SETUP must not require a Go code block, got %q", errStr)
	}
	if !strings.Contains(errStr, "missing func Assert") {
		t.Fatalf("expected missing Assert error, got %q", errStr)
	}
	if !strings.Contains(errStr, "validation errors:") {
		t.Fatalf("expected validation errors header, got %q", errStr)
	}
}

func TestValidationNoErrorForValidTree(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))

	if _, err := DiscoverTreeCases(root); err != nil {
		t.Fatalf("expected no error for valid tree, got %v", err)
	}
}

// Intermediate grouping SETUP may be prose-only (no Go fence); discover and
// hydrate must not emit "must have a Go code block" for that path.
func TestDiscoverAllowsProseOnlyIntermediateSETUP(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{ V int }
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "group/SETUP.md", "# Scenario\n\n**Feature**: organization-only grouping\n\n## Steps\n1. Document only; no Go Setup.\n")
	writeTreeFile(t, root, "group/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { req.V = 1; return nil }
`))
	writeTreeFile(t, root, "group/leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))

	cases, err := DiscoverTreeCases(root)
	if err != nil {
		t.Fatalf("DiscoverTreeCases with prose-only intermediate SETUP: %v", err)
	}
	if len(cases) != 1 {
		t.Fatalf("cases=%d want 1", len(cases))
	}
	light, err := DiscoverTreeCasesLight(root)
	if err != nil {
		t.Fatalf("DiscoverTreeCasesLight: %v", err)
	}
	hydrated, err := HydrateTreeCases(root, light)
	if err != nil {
		t.Fatalf("HydrateTreeCases with prose-only intermediate SETUP: %v", err)
	}
	if len(hydrated) != 1 {
		t.Fatalf("hydrated=%d want 1", len(hydrated))
	}
	// Ensure no setup chain entry forced a Go block error for group/SETUP.md.
	for _, doc := range hydrated[0].SetupFiles {
		if strings.HasSuffix(doc.Path, "group/SETUP.md") || doc.Path == "group/SETUP.md" {
			if doc.GoBlock != nil {
				t.Fatalf("expected nil GoBlock for prose intermediate, got %+v", doc.GoBlock)
			}
		}
	}
}

func TestValidationTestdataDirSkipped(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))
	writeTreeFile(t, root, "testdata/SETUP.md", setupDoc(`
type Request struct{}
type Response struct{}
`))
	writeTreeFile(t, root, "testdata/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "testdata/leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))

	cases, err := DiscoverTreeCases(root)
	if err != nil {
		t.Fatalf("expected no error, testdata/ should be skipped: %v", err)
	}
	if len(cases) != 1 {
		t.Fatalf("expected 1 case (testdata/ skipped), got %d", len(cases))
	}
}

func TestDiscoverTreeCasesSharedSetupMemoized(t *testing.T) {
	// Many leaves share one root SETUP; memoization must not change results.
	root := t.TempDir()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{ X int }
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { req.X = 1; return nil }
`))
	for i := 0; i < 12; i++ {
		name := fmt.Sprintf("leaf%d", i)
		writeTreeFile(t, root, filepath.Join(name, "SETUP.md"), setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }
`))
		writeTreeFile(t, root, filepath.Join(name, "ASSERT.md"), assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))
	}
	cases, err := DiscoverTreeCases(root)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(cases) != 12 {
		t.Fatalf("got %d cases", len(cases))
	}
	for _, c := range cases {
		if len(c.SetupFiles) < 2 {
			t.Fatalf("leaf %s: expected setup chain, got %d", c.Path, len(c.SetupFiles))
		}
	}
}

func TestDiscoverTreeCasesVerbose(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "setup_leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "setup_leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))

	var buf bytes.Buffer
	cases, err := DiscoverTreeCasesVerbose(root, &buf)
	if err != nil {
		t.Fatalf("discover verbose: %v", err)
	}
	if len(cases) != 1 {
		t.Fatalf("expected 1 case, got %d", len(cases))
	}
	out := buf.String()
	if !strings.Contains(out, "DOCTEST.md") {
		t.Fatalf("expected root DOCTEST, got %q", out)
	}
	if !strings.Contains(out, "setup_leaf/SETUP.md") {
		t.Fatalf("expected leaf SETUP, got %q", out)
	}
}

func TestFindModuleRoot(t *testing.T) {
	srcRoot := t.TempDir()
	writeTreeFile(t, srcRoot, "go.mod", "module example.com/a\n\ngo 1.21\n")
	writeTreeFile(t, srcRoot, "a/b/c/SETUP.md", "# test\n")

	modRoot, modPath, ok := FindModuleRoot(filepath.Join(srcRoot, "a/b/c"))
	if !ok {
		t.Fatal("expected to find module root")
	}
	if modRoot != srcRoot {
		t.Fatalf("expected mod root %q, got %q", srcRoot, modRoot)
	}
	if modPath != "example.com/a" {
		t.Fatalf("expected mod path example.com/a, got %q", modPath)
	}
}

func TestFindModuleRootNotFound(t *testing.T) {
	root := t.TempDir()
	_, _, ok := FindModuleRoot(root)
	if ok {
		t.Fatal("expected no module root")
	}
}

func TestAssembleTestSourceIncludesDoctestSessionID(t *testing.T) {
	root := t.TempDir()
	tc := TreeCase{
		Name: "leaf",
		Path: "leaf",
		SetupFiles: []SetupDocument{{
			Path: "",
			GoBlock: &GoBlock{
				Run: &FuncSnippet{
					Name:    "Run",
					Params:  "t *testing.T, d *session.Doctest, req *Request",
					Results: "(*Response, error)",
					Body:    "{ return &Response{}, nil }",
				},
				Types: map[string]bool{"Request": true, "Response": true},
				TypeDecls: []string{
					"type Request struct{}",
					"type Response struct{}",
				},
			},
		}},
		AssertFile: AssertDocument{
			GoBlock: GoBlock{
				Assert: &FuncSnippet{
					Name:   "Assert",
					Params: "t *testing.T, d *session.Doctest, req *Request, resp *Response, err error",
					Body:   "{}",
				},
			},
		},
	}

	src, err := AssembleTestSource(tc, false, "leaf_tc", root)
	if err != nil {
		t.Fatalf("assemble: %v", err)
	}
	assertAssembledMatchesFixture(t, src, root, "assemble_session_id.go.fixture")
}

func TestWriteGoModSkipsWhenSourceUnchanged(t *testing.T) {
	modRoot := t.TempDir()
	genDir := filepath.Join(t.TempDir(), "gen")
	writeTreeFile(t, modRoot, "go.mod", "module example.com/a\n\ngo 1.21\n")

	if err := WriteGoMod(genDir, modRoot, "example.com/a", true, false, "", false, ""); err != nil {
		t.Fatalf("first WriteGoMod: %v", err)
	}
	first, err := os.ReadFile(filepath.Join(genDir, "go.mod"))
	if err != nil {
		t.Fatalf("read first go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(genDir, "doctest.tidy-done"), []byte("done"), 0644); err != nil {
		t.Fatalf("write tidy marker: %v", err)
	}

	if err := WriteGoMod(genDir, modRoot, "example.com/a", true, false, "", false, ""); err != nil {
		t.Fatalf("second WriteGoMod: %v", err)
	}
	second, err := os.ReadFile(filepath.Join(genDir, "go.mod"))
	if err != nil {
		t.Fatalf("read second go.mod: %v", err)
	}
	if string(first) != string(second) {
		t.Fatalf("expected cached go.mod to be unchanged")
	}
	if _, err := os.Stat(filepath.Join(genDir, "doctest.tidy-done")); err != nil {
		t.Fatalf("expected tidy marker to remain when source go.mod unchanged: %v", err)
	}
}

func TestWriteGoModRegeneratesWhenSourceChanges(t *testing.T) {
	modRoot := t.TempDir()
	genDir := filepath.Join(t.TempDir(), "gen")
	writeTreeFile(t, modRoot, "go.mod", "module example.com/a\n\ngo 1.21\n\nreplace localdep => ./dep\n")

	if err := WriteGoMod(genDir, modRoot, "example.com/a", true, false, "", false, ""); err != nil {
		t.Fatalf("first WriteGoMod: %v", err)
	}
	if !strings.Contains(readFileString(t, filepath.Join(genDir, "go.mod")), "replace localdep") {
		t.Fatal("expected first go.mod to include replace localdep")
	}
	if err := os.WriteFile(filepath.Join(genDir, "doctest.tidy-done"), []byte("done"), 0644); err != nil {
		t.Fatalf("write tidy marker: %v", err)
	}

	writeTreeFile(t, modRoot, "go.mod", "module example.com/a\n\ngo 1.21\n")
	if err := WriteGoMod(genDir, modRoot, "example.com/a", true, false, "", false, ""); err != nil {
		t.Fatalf("second WriteGoMod: %v", err)
	}
	updated := readFileString(t, filepath.Join(genDir, "go.mod"))
	if strings.Contains(updated, "replace localdep") {
		t.Fatalf("expected regenerated go.mod to drop stale replace:\n%s", updated)
	}
	if _, err := os.Stat(filepath.Join(genDir, "doctest.tidy-done")); !os.IsNotExist(err) {
		t.Fatalf("expected tidy marker removed after source go.mod change, err=%v", err)
	}
}

func TestWriteGoModSkipsAssertSubmoduleReplaceForDoctestModule(t *testing.T) {
	modRoot := t.TempDir()
	genDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(modRoot, "go.mod"), []byte("module github.com/xhd2015/doctest\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := WriteGoMod(genDir, modRoot, "github.com/xhd2015/doctest", true, true, "/tmp/assert-cache", false, ""); err != nil {
		t.Fatalf("WriteGoMod: %v", err)
	}
	goMod := readFileString(t, filepath.Join(genDir, "go.mod"))
	if !strings.Contains(goMod, "replace github.com/xhd2015/doctest => "+modRoot) {
		t.Fatalf("expected parent module replace, got:\n%s", goMod)
	}
	if strings.Contains(goMod, "replace github.com/xhd2015/doctest/assert =>") {
		t.Fatalf("expected no assert submodule replace for doctest self-tests, got:\n%s", goMod)
	}
}

// Multi-tree ./... prepare calls WriteGoMod with different per-tree assert flags.
// For the doctest module those flags do not affect go.mod content; fingerprint
// and content-stable writes must not rewrite go.mod or drop tidy-done (mtime
// of gen root is a go testcache input via package-dir chdir/stat).
func TestWriteGoModDoctestModuleIgnoresIneffectiveAssertFlags(t *testing.T) {
	modRoot := t.TempDir()
	genDir := filepath.Join(t.TempDir(), "gen")
	if err := os.WriteFile(filepath.Join(modRoot, "go.mod"), []byte("module github.com/xhd2015/doctest\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := WriteGoMod(genDir, modRoot, "github.com/xhd2015/doctest", true, false, "", false, ""); err != nil {
		t.Fatalf("first WriteGoMod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(genDir, "doctest.tidy-done"), []byte("done"), 0644); err != nil {
		t.Fatalf("write tidy marker: %v", err)
	}
	modPath := filepath.Join(genDir, "go.mod")
	before, err := os.Stat(modPath)
	if err != nil {
		t.Fatal(err)
	}

	// Second call as if another tree needs assert import — ineffective for doctest.
	if err := WriteGoMod(genDir, modRoot, "github.com/xhd2015/doctest", true, true, "/tmp/assert-cache-"+t.Name(), true, "/tmp/session-cache"); err != nil {
		t.Fatalf("second WriteGoMod: %v", err)
	}
	after, err := os.Stat(modPath)
	if err != nil {
		t.Fatal(err)
	}
	if !before.ModTime().Equal(after.ModTime()) {
		t.Fatalf("go.mod mtime changed under ineffective assert/session flags: %v -> %v", before.ModTime(), after.ModTime())
	}
	if _, err := os.Stat(filepath.Join(genDir, "doctest.tidy-done")); err != nil {
		t.Fatalf("expected tidy marker retained: %v", err)
	}
}

func readFileString(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func TestParseHelperSharedTypeNamedResultsSetsResultTypes(t *testing.T) {
	block := GoBlock{SourcePath: "DOCTEST.md", Code: `
import "testing"
func pickTwoPorts(base int) (port, alt int) { return base, base + 1 }
`}
	if err := parseGoBlock(&block); err != nil {
		t.Fatalf("parseGoBlock: %v", err)
	}
	if len(block.Helpers) != 1 {
		t.Fatalf("expected 1 helper, got %d", len(block.Helpers))
	}
	h := block.Helpers[0]
	if h.Results != "port int, alt int" {
		t.Fatalf("Results = %q", h.Results)
	}
	if h.ResultTypes != "(int, int)" {
		t.Fatalf("ResultTypes = %q, want (int, int)", h.ResultTypes)
	}
}

func TestResultTypesStringStripsNamedReturns(t *testing.T) {
	code := `package p
func startTestServer(t *testing.T) (base string, cleanup func()) { return "", nil }
func pickFreePort(base int) (int, error) { return 0, nil }
func pickTwoPorts(base int) (port, alt int) { return 0, 0 }`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "helpers.go", code, 0)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var gotStart, gotPort, gotTwo string
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		switch fn.Name.Name {
		case "startTestServer":
			gotStart = resultTypesString(fset, fn.Type.Results)
		case "pickFreePort":
			gotPort = resultTypesString(fset, fn.Type.Results)
		case "pickTwoPorts":
			gotTwo = resultTypesString(fset, fn.Type.Results)
		}
	}
	if gotStart != "(string, func())" {
		t.Fatalf("startTestServer result types = %q, want (string, func())", gotStart)
	}
	if gotPort != "(int, error)" {
		t.Fatalf("pickFreePort result types = %q, want (int, error)", gotPort)
	}
	if gotTwo != "(int, int)" {
		t.Fatalf("pickTwoPorts result types = %q, want (int, int)", gotTwo)
	}
}

func TestAssembleFuncClosureSharedTypeNamedResultsParses(t *testing.T) {
	root := t.TempDir()
	tc := TreeCase{
		Name: "leaf",
		Path: "leaf",
		SetupFiles: []SetupDocument{{
			Path: "",
			GoBlock: &GoBlock{
				Run: &FuncSnippet{
					Name:        "Run",
					Params:      "t *testing.T, d *session.Doctest, req *Request",
					Results:     "(*Response, error)",
					ResultTypes: "(*Response, error)",
					Body:        "{ return &Response{}, nil }",
				},
				Helpers: []FuncSnippet{{
					Name:           "pickTwoPorts",
					Params:         "base int",
					Results:        "port int, alt int",
					ResultTypes:    "(int, int)",
					ClosureResults: "(port int, alt int)",
					Body:           "{ return base, base + 1 }",
				}},
				Types: map[string]bool{"Request": true, "Response": true},
				TypeDecls: []string{
					"type Request struct{}",
					"type Response struct{}",
				},
			},
		}},
		AssertFile: AssertDocument{
			GoBlock: GoBlock{
				Assert: &FuncSnippet{
					Name:   "Assert",
					Params: "t *testing.T, d *session.Doctest, req *Request, resp *Response, err error",
					Body:   "{}",
				},
			},
		},
	}

	src, err := AssembleTestSource(tc, false, "leaf_tc", root)
	if err != nil {
		t.Fatalf("assemble: %v", err)
	}
	// Full assembled source is the reviewable contract (top-level helpers + named results).
	assertAssembledMatchesFixture(t, src, root, "assemble_named_results_pick_two_ports.go.fixture")
	if _, err := parser.ParseFile(token.NewFileSet(), "generated_test.go", src, 0); err != nil {
		t.Fatalf("generated source should parse: %v\n%s", err, src)
	}
	if _, err := formatGeneratedGo("generated_test.go", []byte(src)); err != nil {
		t.Fatalf("formatGeneratedGo should succeed: %v\n%s", err, src)
	}
}

func parseHelpersBlock(t *testing.T, code string) *GoBlock {
	t.Helper()
	block := GoBlock{SourcePath: "SETUP.md", Code: code}
	if err := parseGoBlock(&block); err != nil {
		t.Fatalf("parseGoBlock: %v", err)
	}
	return &block
}

// TestAssembleFuncClosurePreservesNamedResults reproduces a helper whose body
// assigns to named result variables (like wrk's setupWrkWorktreeFromMain).
// Stripping names to type-only signatures leaves those assignments referencing
// undeclared identifiers and breaks compilation. The closure signature must
// keep the names, parenthesized.
func TestAssembleFuncClosurePreservesNamedResults(t *testing.T) {
	root := t.TempDir()
	block := parseHelpersBlock(t, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
func splitNames(req *Request) (mainRepo, wtDir, branch string) {
	mainRepo = "a"
	wtDir = "b"
	branch = "c"
	return
}
`)
	tc := TreeCase{
		Name: "leaf",
		Path: "leaf",
		SetupFiles: []SetupDocument{{Path: "", GoBlock: block}},
		AssertFile: AssertDocument{GoBlock: GoBlock{
			Assert: &FuncSnippet{
				Name:   "Assert",
				Params: "t *testing.T, d *session.Doctest, req *Request, resp *Response, err error",
				Body:   "{}",
			},
		}},
	}

	src, err := AssembleTestSource(tc, false, "leaf_tc", root)
	if err != nil {
		t.Fatalf("assemble: %v", err)
	}
	// Named results kept on top-level helper (not type-only strip). Fixture is the review surface.
	assertAssembledMatchesFixture(t, src, root, "assemble_named_results_split_names.go.fixture")
	if strings.Contains(src, "func splitNames(req *Request) (string, string, string)") {
		t.Fatalf("named results should not be stripped to type-only, got:\n%s", src)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "generated_test.go", src, 0); err != nil {
		t.Fatalf("generated source should parse: %v\n%s", err, src)
	}
	if _, err := formatGeneratedGo("generated_test.go", []byte(src)); err != nil {
		t.Fatalf("formatGeneratedGo should succeed: %v\n%s", err, src)
	}
}

// TestAssembleHelpersTopoSortedForwardRef reproduces two helpers where the
// caller is defined before the callee (legal for top-level funcs, illegal for
// closures). The codegen must emit the callee before the caller.
func TestAssembleHelpersTopoSortedForwardRef(t *testing.T) {
	root := t.TempDir()
	block := parseHelpersBlock(t, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
func caller(x int) string { return callee(x) }
func callee(x int) string { return "x" }
`)
	tc := TreeCase{
		Name: "leaf",
		Path: "leaf",
		SetupFiles: []SetupDocument{{Path: "", GoBlock: block}},
		AssertFile: AssertDocument{GoBlock: GoBlock{
			Assert: &FuncSnippet{
				Name:   "Assert",
				Params: "t *testing.T, d *session.Doctest, req *Request, resp *Response, err error",
				Body:   "{}",
			},
		}},
	}

	src, err := AssembleTestSource(tc, false, "leaf_tc", root)
	if err != nil {
		t.Fatalf("assemble: %v", err)
	}
	// Fixture shows callee before caller (topo order) as top-level funcs.
	assertAssembledMatchesFixture(t, src, root, "assemble_helpers_topo_sorted.go.fixture")
	calleeIdx := strings.Index(src, "func callee(")
	callerIdx := strings.Index(src, "func caller(")
	if calleeIdx < 0 || callerIdx < 0 {
		t.Fatalf("expected both helper funcs in source, got:\n%s", src)
	}
	if calleeIdx >= callerIdx {
		t.Fatalf("expected callee declared before caller (calleeIdx=%d callerIdx=%d), got:\n%s", calleeIdx, callerIdx, src)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "generated_test.go", src, 0); err != nil {
		t.Fatalf("generated source should parse: %v\n%s", err, src)
	}
	if _, err := formatGeneratedGo("generated_test.go", []byte(src)); err != nil {
		t.Fatalf("formatGeneratedGo should succeed: %v\n%s", err, src)
	}
}

// assertAssembledMatchesFixture compares AssembleTestSource output to a golden
// fixture under testdata/. Absolute roots are rewritten to {{DOCTEST_ROOT}}.
// Set UPDATE_FIXTURES=1 to regenerate fixtures for review.
func assertAssembledMatchesFixture(t *testing.T, got, root, fixtureName string) {
	t.Helper()
	absRoot, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("abs root: %v", err)
	}
	normalized := strings.ReplaceAll(got, absRoot, "{{DOCTEST_ROOT}}")
	normalized = strings.ReplaceAll(normalized, root, "{{DOCTEST_ROOT}}")
	if !strings.HasSuffix(normalized, "\n") {
		normalized += "\n"
	}

	fixturePath := filepath.Join("testdata", fixtureName)
	if os.Getenv("UPDATE_FIXTURES") == "1" {
		if err := os.MkdirAll(filepath.Dir(fixturePath), 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}
		if err := os.WriteFile(fixturePath, []byte(normalized), 0o644); err != nil {
			t.Fatalf("write fixture %s: %v", fixturePath, err)
		}
		t.Logf("updated fixture %s", fixturePath)
		return
	}

	wantBytes, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture %s: %v (run with UPDATE_FIXTURES=1 to create)", fixturePath, err)
	}
	if normalized != string(wantBytes) {
		t.Fatalf("assembled source does not match fixture %s\n--- got ---\n%s\n--- want ---\n%s", fixtureName, normalized, string(wantBytes))
	}
}
