package core

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootSetupRequiresRequestAndResponseTypes(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
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
	if !strings.Contains(err.Error(), "Request") || !strings.Contains(err.Error(), "Response") {
		t.Fatalf("expected Request/Response error, got %v", err)
	}
}

func TestValidationReportsAllErrorsAtOnce(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "SETUP.md", "# Setup\n\n```go\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n```\n")
	writeTreeFile(t, root, "leaf/SETUP.md", "# Setup\n\nProse only, no Go block.\n")
	writeTreeFile(t, root, "leaf/ASSERT.md", "# Assert\n\n```go\nfunc Check(t *testing.T, req *Request, resp *Response, err error) {}\n```\n")

	_, err := DiscoverTreeCases(root)
	if err == nil {
		t.Fatal("expected validation error")
	}
	errStr := err.Error()
	if !strings.Contains(errStr, "must define type Request") {
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
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
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
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
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
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
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
	if !strings.Contains(out, "SETUP.md") {
		t.Fatalf("expected root SETUP, got %q", out)
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
