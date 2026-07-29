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

// createPathScopeMidNestedRoot: mid + sibling + nested DOCTEST under mid (same module).
// Used for S1: ./tree/mid/... must be one go test ./tree/mid/... covering nested.
func createPathScopeMidNestedRoot(t *testing.T) string {
	t.Helper()
	proj := createPathScopeMidSibling(t)
	tree := filepath.Join(proj, "tree")
	nested := filepath.Join(tree, "mid", "nested")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(pathScopeMinimalRun())), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "SETUP.md"), []byte(pathScopeGoBlock(
		"import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }",
	)), 0644); err != nil {
		t.Fatal(err)
	}
	pathScopeWriteLeaf(t, filepath.Join(nested, "three"), "NESTED_LEAF")
	return proj
}

// createPathScopeMidNestedGomod: mid + sibling + nested go.mod under mid with DOCTEST.
// Used for S4: same-gen must combine patterns; multi-gen must use different cd.
func createPathScopeMidNestedGomod(t *testing.T) string {
	t.Helper()
	proj := createPathScopeMidSibling(t)
	// rewrite module path for nested child path policy
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module midtreeproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	nestedMod := filepath.Join(proj, "tree", "mid", "nestedmod")
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
	if err := os.WriteFile(filepath.Join(suite, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(pathScopeMinimalRun())), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(suite, "SETUP.md"), []byte(pathScopeGoBlock(
		"import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }",
	)), 0644); err != nil {
		t.Fatal(err)
	}
	pathScopeWriteLeaf(t, filepath.Join(suite, "one"), "NESTED_MOD_LEAF")
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

// pathScopeCdDir extracts the directory from "cd DIR && go test ..." plan lines.
// Returns "" if the line has no cd prefix.
func pathScopeCdDir(plan string) string {
	plan = strings.TrimSpace(plan)
	if !strings.HasPrefix(plan, "cd ") {
		return ""
	}
	rest := strings.TrimPrefix(plan, "cd ")
	i := strings.Index(rest, " && ")
	if i < 0 {
		return rest
	}
	return strings.TrimSpace(rest[:i])
}

// pathScopeAssertPlansByDir: same RunDir must appear at most once (combine patterns).
func pathScopeAssertPlansByDir(t *testing.T, plans []string, out string) {
	t.Helper()
	if len(plans) == 0 {
		t.Fatalf("no go test plan lines:\n%s", out)
	}
	byDir := map[string][]string{}
	for _, p := range plans {
		d := pathScopeCdDir(p)
		if d == "" {
			d = "."
		}
		byDir[d] = append(byDir[d], p)
	}
	for d, ps := range byDir {
		if len(ps) > 1 {
			t.Fatalf("same cd %q has %d go test invocations (must combine into one go test):\n  %s\nfull:\n%s",
				d, len(ps), strings.Join(ps, "\n  "), out)
		}
	}
}

func pathScopeCountMarker(out, marker string) int {
	return strings.Count(out, "MARKER:"+marker)
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout < 2*time.Minute {
		req.Timeout = 2 * time.Minute
	}
	return nil
}
```
