package core

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootSetupRequiresRequestAndResponseTypes(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
func Run(t *testing.T, req *Request) (*Response, error) { return nil, nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
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
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`func Run(t *testing.T, req *Request) (*Response, error) { return nil, nil }`))
	writeTreeFile(t, root, "SETUP.md", "# Setup\n\n```go\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n```\n")
	writeTreeFile(t, root, "leaf/SETUP.md", "# Setup\n\nProse only, no Go block.\n")
	writeTreeFile(t, root, "leaf/ASSERT.md", "# Assert\n\n```go\nfunc Check(t *testing.T, req *Request, resp *Response, err error) {}\n```\n")

	_, err := DiscoverTreeCases(root)
	if err == nil {
		t.Fatal("expected validation error")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "must define type Request") && !strings.Contains(errStr, "type Response") {
		t.Fatalf("expected missing types error, got %q", errStr)
	}
	if !strings.Contains(errStr, "must have a Go code block") {
		t.Fatalf("expected missing Go block error, got %q", errStr)
	}
	if !strings.Contains(errStr, "missing func Assert") {
		t.Fatalf("expected missing Assert error, got %q", errStr)
	}
	if !strings.Contains(errStr, "validation errors:") {
		t.Fatalf("expected validation errors header, got %q", errStr)
	}
	lines := strings.Split(strings.TrimSpace(errStr), "\n")
	if len(lines) < 4 {
		t.Fatalf("expected at least 3 errors, got %q", errStr)
	}
}

func TestValidationNoErrorForValidTree(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	if _, err := DiscoverTreeCases(root); err != nil {
		t.Fatalf("expected no error for valid tree, got %v", err)
	}
}

func TestValidationTestdataDirSkipped(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))
	writeTreeFile(t, root, "testdata/SETUP.md", setupDoc(`
type Request struct{}
type Response struct{}
`))
	writeTreeFile(t, root, "testdata/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "testdata/leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	cases, err := DiscoverTreeCases(root)
	if err != nil {
		t.Fatalf("expected no error, testdata/ should be skipped: %v", err)
	}
	if len(cases) != 1 {
		t.Fatalf("expected 1 case (testdata/ skipped), got %d", len(cases))
	}
}

func TestDiscoverTreeCasesVerbose(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "setup_leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "setup_leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
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
					Params:  "t *testing.T, req *Request",
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
					Params: "t *testing.T, req *Request, resp *Response, err error",
					Body:   "{}",
				},
			},
		},
	}

	src, err := AssembleTestSource(tc, false, "leaf_tc", root)
	if err != nil {
		t.Fatalf("assemble: %v", err)
	}
	if !strings.Contains(src, "DOCTEST_SESSION_ID, __sessionOk := syscall.Getenv(\"DOCTEST_SESSION_ID\")") {
		t.Fatalf("expected DOCTEST_SESSION_ID assignment, got:\n%s", src)
	}
	if !strings.Contains(src, "t.Fatalf(\"DOCTEST_SESSION_ID not set\")") {
		t.Fatalf("expected DOCTEST_SESSION_ID missing fatal, got:\n%s", src)
	}
	if !strings.Contains(src, "\"syscall\"") {
		t.Fatalf("expected syscall import, got:\n%s", src)
	}
}

func TestWriteGoModSkipsWhenSourceUnchanged(t *testing.T) {
	modRoot := t.TempDir()
	genDir := filepath.Join(t.TempDir(), "gen")
	writeTreeFile(t, modRoot, "go.mod", "module example.com/a\n\ngo 1.21\n")

	if err := WriteGoMod(genDir, modRoot, "example.com/a", true, false, ""); err != nil {
		t.Fatalf("first WriteGoMod: %v", err)
	}
	first, err := os.ReadFile(filepath.Join(genDir, "go.mod"))
	if err != nil {
		t.Fatalf("read first go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(genDir, "doctest.tidy-done"), []byte("done"), 0644); err != nil {
		t.Fatalf("write tidy marker: %v", err)
	}

	if err := WriteGoMod(genDir, modRoot, "example.com/a", true, false, ""); err != nil {
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

	if err := WriteGoMod(genDir, modRoot, "example.com/a", true, false, ""); err != nil {
		t.Fatalf("first WriteGoMod: %v", err)
	}
	if !strings.Contains(readFileString(t, filepath.Join(genDir, "go.mod")), "replace localdep") {
		t.Fatal("expected first go.mod to include replace localdep")
	}
	if err := os.WriteFile(filepath.Join(genDir, "doctest.tidy-done"), []byte("done"), 0644); err != nil {
		t.Fatalf("write tidy marker: %v", err)
	}

	writeTreeFile(t, modRoot, "go.mod", "module example.com/a\n\ngo 1.21\n")
	if err := WriteGoMod(genDir, modRoot, "example.com/a", true, false, ""); err != nil {
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

func readFileString(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
