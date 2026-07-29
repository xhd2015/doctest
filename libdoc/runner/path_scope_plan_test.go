package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Path-scoped go test plan shape (S1 / S4) — classic TDD locks.
//
// Rule: group by RunDir (cd). Same Dir → one go test (combine/collapse).
// Different Dir → one go test per Dir. Never two processes for the same cd.

func extractGoTestPlans(out string) []string {
	var plans []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "&& go test") {
			plans = append(plans, line)
		}
	}
	return plans
}

func planCdDir(plan string) string {
	plan = strings.TrimSpace(plan)
	if !strings.HasPrefix(plan, "cd ") {
		return "."
	}
	rest := strings.TrimPrefix(plan, "cd ")
	if i := strings.Index(rest, " && "); i >= 0 {
		return strings.TrimSpace(rest[:i])
	}
	return rest
}

func assertPlansSameDirAtMostOnce(t *testing.T, plans []string, out string) {
	t.Helper()
	byDir := map[string][]string{}
	for _, p := range plans {
		d := planCdDir(p)
		byDir[d] = append(byDir[d], p)
	}
	for d, ps := range byDir {
		if len(ps) > 1 {
			t.Fatalf("same cd %q has %d go test cmds (must combine):\n  %s\nfull:\n%s",
				d, len(ps), strings.Join(ps, "\n  "), out)
		}
	}
}

func pathScopeMarkerCount(out, name string) int {
	return strings.Count(out, "MARKER:"+name)
}

func writePathScopeLeaf(t *testing.T, dir, marker string) {
	t.Helper()
	writeTreeFile(t, dir, "SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Logf("MARKER:`+marker+`")
	return nil
}
`))
	writeTreeFile(t, dir, "ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))
}

// fixture: mid + sibling + nested DOCTEST under mid (same module).
func fixtureMidNestedRoot(t *testing.T) string {
	t.Helper()
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module pathscope\ngo 1.21\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	tree := filepath.Join(proj, "tree")
	writeTreeFile(t, tree, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, tree, "SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }
`))
	writePathScopeLeaf(t, filepath.Join(tree, "mid", "two"), "MID_LEAF")
	writePathScopeLeaf(t, filepath.Join(tree, "sibling", "one"), "SIBLING_LEAF")
	nested := filepath.Join(tree, "mid", "nested")
	writeTreeFile(t, nested, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, nested, "SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }
`))
	writePathScopeLeaf(t, filepath.Join(nested, "three"), "NESTED_LEAF")
	return proj
}

// fixture: mid + sibling + nested go.mod under mid.
func fixtureMidNestedGomod(t *testing.T) string {
	t.Helper()
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module midtreeproj\ngo 1.21\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	tree := filepath.Join(proj, "tree")
	writeTreeFile(t, tree, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, tree, "SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }
`))
	writePathScopeLeaf(t, filepath.Join(tree, "mid", "two"), "MID_LEAF")
	writePathScopeLeaf(t, filepath.Join(tree, "sibling", "one"), "SIBLING_LEAF")
	nestedMod := filepath.Join(tree, "mid", "nestedmod")
	if err := os.MkdirAll(nestedMod, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nestedMod, "go.mod"), []byte("module midtreeproj/nestedmod\ngo 1.21\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	suite := filepath.Join(nestedMod, "suite")
	writeTreeFile(t, suite, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, suite, "SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }
`))
	writePathScopeLeaf(t, filepath.Join(suite, "one"), "NESTED_MOD_LEAF")
	return proj
}

func runPathScopedMid(t *testing.T, proj string) (out string, err error) {
	t.Helper()
	gen := filepath.Join(t.TempDir(), "gen")
	if err := os.MkdirAll(gen, 0o755); err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(proj); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	var stdout, stderr bytes.Buffer
	runErr := TestWithWriters([]string{
		"-v", "--no-color", "--label-all", "--gen-dir", gen, "-count=1", "./tree/mid/...",
	}, &stdout, &stderr)
	return stdout.String() + "\n" + stderr.String(), runErr
}

// TestPathScoped_S1_MidNestedRoot_SingleGoTestPlan: parent ./tree/mid/... covers
// nested DOCTEST packages — one go test, nested marker once.
func TestPathScoped_S1_MidNestedRoot_SingleGoTestPlan(t *testing.T) {
	proj := fixtureMidNestedRoot(t)
	out, runErr := runPathScopedMid(t, proj)
	if runErr != nil {
		// Still assert plan shape on partial output when useful.
		t.Logf("run err (plan asserts still apply): %v", runErr)
	}
	plans := extractGoTestPlans(out)
	assertPlansSameDirAtMostOnce(t, plans, out)
	if len(plans) != 1 {
		t.Fatalf("S1 want exactly 1 go test plan, got %d:\n  %s\nfull:\n%s",
			len(plans), strings.Join(plans, "\n  "), out)
	}
	if !strings.Contains(plans[0], "./tree/mid/...") {
		t.Fatalf("plan want ./tree/mid/...:\n  %s\n%s", plans[0], out)
	}
	if n := pathScopeMarkerCount(out, "MID_LEAF"); n != 1 {
		t.Fatalf("MARKER:MID_LEAF want 1 got %d\n%s", n, out)
	}
	if n := pathScopeMarkerCount(out, "NESTED_LEAF"); n != 1 {
		t.Fatalf("MARKER:NESTED_LEAF want 1 got %d\n%s", n, out)
	}
	if pathScopeMarkerCount(out, "SIBLING_LEAF") != 0 {
		t.Fatalf("sibling must not run\n%s", out)
	}
	if runErr != nil {
		t.Fatalf("run failed: %v\n%s", runErr, out)
	}
}

// TestPathScoped_S4_MidNestedGomod_SameGen_CombinedPlan: shared gen-dir → one
// go test with both patterns (or multi-gen with different cd — not same cd twice).
func TestPathScoped_S4_MidNestedGomod_SameGen_CombinedPlan(t *testing.T) {
	proj := fixtureMidNestedGomod(t)
	out, runErr := runPathScopedMid(t, proj)
	if runErr != nil {
		t.Logf("run err (plan asserts still apply): %v", runErr)
	}
	plans := extractGoTestPlans(out)
	assertPlansSameDirAtMostOnce(t, plans, out)
	if len(plans) != 1 {
		t.Fatalf("S4 shared-gen want exactly 1 combined go test, got %d:\n  %s\nfull:\n%s",
			len(plans), strings.Join(plans, "\n  "), out)
	}
	p := plans[0]
	hasSuite := strings.Contains(p, "./suite") || strings.Contains(p, "suite/...")
	hasMid := strings.Contains(p, "mid") && strings.Contains(p, "/...")
	if !hasSuite || !hasMid {
		t.Fatalf("combined plan want suite + mid patterns:\n  %s\nfull:\n%s", p, out)
	}
	if n := pathScopeMarkerCount(out, "MID_LEAF"); n != 1 {
		t.Fatalf("MARKER:MID_LEAF want 1 got %d\n%s", n, out)
	}
	if n := pathScopeMarkerCount(out, "NESTED_MOD_LEAF"); n != 1 {
		t.Fatalf("MARKER:NESTED_MOD_LEAF want 1 got %d\n%s", n, out)
	}
	if pathScopeMarkerCount(out, "SIBLING_LEAF") != 0 {
		t.Fatalf("sibling must not run\n%s", out)
	}
	if runErr != nil {
		t.Fatalf("run failed: %v\n%s", runErr, out)
	}
}
