# Scenario

**Feature**: `./tree/mid/...` — no siblings; nested DOCTEST / nested go.mod+DOCTEST under prefix

```
tree/                 DOCTEST
  sibling/one/        SIBLING_LEAF     ✗
  mid/two/            MID_LEAF         ✓
  mid/nested/         DOCTEST          ✓  (same module)
  mid/nestedmod/      go.mod + DOCTEST ✓  (child module path; policy B/C)
```

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

var midTreeBt = "`" + "`" + "`"

func midTreeGoBlock(code string) string {
	return "## Test\n\n" + midTreeBt + "go\n" + code + "\n" + midTreeBt + "\n"
}

func midTreeMinimalRun() string {
	return `import "testing"
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }`
}

func midTreeLeaf(marker string) (setup, assert string) {
	setup = midTreeGoBlock("import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error {\n\tt.Logf(\"MARKER:" + marker + "\")\n\treturn nil\n}")
	assert = midTreeGoBlock("import \"testing\"\nfunc Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}")
	return setup, assert
}

func midTreeWriteLeaf(t *testing.T, dir, marker string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	s, a := midTreeLeaf(marker)
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(s), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ASSERT.md"), []byte(a), 0644); err != nil {
		t.Fatal(err)
	}
}

// createMidTreePrefixProject: sibling + mid leaf + nested DOCTEST under mid.
func createMidTreePrefixProject(t *testing.T) string {
	t.Helper()
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module midtreeproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	tree := filepath.Join(proj, "tree")
	if err := os.MkdirAll(tree, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tree, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(midTreeMinimalRun())), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tree, "SETUP.md"), []byte(midTreeGoBlock(
		"import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }",
	)), 0644); err != nil {
		t.Fatal(err)
	}
	midTreeWriteLeaf(t, filepath.Join(tree, "sibling", "one"), "SIBLING_LEAF")
	midTreeWriteLeaf(t, filepath.Join(tree, "mid", "two"), "MID_LEAF")

	nested := filepath.Join(tree, "mid", "nested")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(midTreeMinimalRun())), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "SETUP.md"), []byte(midTreeGoBlock(
		"import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }",
	)), 0644); err != nil {
		t.Fatal(err)
	}
	midTreeWriteLeaf(t, filepath.Join(nested, "three"), "NESTED_LEAF")
	return proj
}

// createMidTreeNestedGomodProject: sibling + mid leaf + nested go.mod (child path) with DOCTEST.
// Nested module path midtreeproj/nestedmod is a child of midtreeproj (policy C discoverability).
func createMidTreeNestedGomodProject(t *testing.T) string {
	t.Helper()
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module midtreeproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	tree := filepath.Join(proj, "tree")
	if err := os.MkdirAll(tree, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tree, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(midTreeMinimalRun())), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tree, "SETUP.md"), []byte(midTreeGoBlock(
		"import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }",
	)), 0644); err != nil {
		t.Fatal(err)
	}
	midTreeWriteLeaf(t, filepath.Join(tree, "sibling", "one"), "SIBLING_LEAF")
	midTreeWriteLeaf(t, filepath.Join(tree, "mid", "two"), "MID_LEAF")

	nestedMod := filepath.Join(tree, "mid", "nestedmod")
	if err := os.MkdirAll(nestedMod, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nestedMod, "go.mod"), []byte("module midtreeproj/nestedmod\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	suite := filepath.Join(nestedMod, "suite")
	if err := os.MkdirAll(suite, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(suite, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(midTreeMinimalRun())), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(suite, "SETUP.md"), []byte(midTreeGoBlock(
		"import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }",
	)), 0644); err != nil {
		t.Fatal(err)
	}
	midTreeWriteLeaf(t, filepath.Join(suite, "one"), "NESTED_MOD_LEAF")
	return proj
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout < 2*time.Minute {
		req.Timeout = 2 * time.Minute
	}
	return nil
}
```
