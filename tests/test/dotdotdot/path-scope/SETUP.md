# Scenario

**Feature**: path scope for **both** phases — generate and go test

User path argument (with or without `...`) is a hard boundary:

| Phase | Rule |
|--------|------|
| **Generate** | Only create/update gen content under the selected subpath. Content **outside** that path must not be modified. |
| **go test** | Only run packages/cases under that same subpath. Mid vs sibling must not share one root suite plan. |

Related: `path-prefix/mid-tree-dotdotdot` (execution markers), `gen-manifest-scope` (multi-tree ledger). This group locks **in-tree mid vs sibling** for both phases.

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

var pathScopeBt = "`" + "`" + "`"

func pathScopeGoBlock(code string) string {
	return "## Test\n\n" + pathScopeBt + "go\n" + code + "\n" + pathScopeBt + "\n"
}

func pathScopeMinimalRun() string {
	return `import "testing"
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }`
}

func pathScopeWriteLeaf(t *testing.T, dir, marker string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	setup := pathScopeGoBlock("import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error {\n\tt.Logf(\"MARKER:" + marker + "\")\n\treturn nil\n}")
	assert := pathScopeGoBlock("import \"testing\"\nfunc Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}")
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(setup), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ASSERT.md"), []byte(assert), 0644); err != nil {
		t.Fatal(err)
	}
}

// createPathScopeMidSibling: one module, one DOCTEST tree, mid + sibling leaves.
func createPathScopeMidSibling(t *testing.T) string {
	t.Helper()
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module pathscopeproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	tree := filepath.Join(proj, "tree")
	if err := os.MkdirAll(tree, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tree, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(pathScopeMinimalRun())), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tree, "SETUP.md"), []byte(pathScopeGoBlock(
		"import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }",
	)), 0644); err != nil {
		t.Fatal(err)
	}
	pathScopeWriteLeaf(t, filepath.Join(tree, "mid", "two"), "MID_LEAF")
	pathScopeWriteLeaf(t, filepath.Join(tree, "sibling", "one"), "SIBLING_LEAF")
	return proj
}

func pathScopeOut(resp *Response) string {
	if resp == nil {
		return ""
	}
	return resp.Stdout + "\n" + resp.Stderr
}

func pathScopeGoTestPlanLines(out string) []string {
	var plans []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "&& go test") {
			plans = append(plans, line)
		}
	}
	return plans
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout < 2*time.Minute {
		req.Timeout = 2 * time.Minute
	}
	return nil
}
```
