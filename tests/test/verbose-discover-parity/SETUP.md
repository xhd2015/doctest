# Scenario

**Feature**: lock `-v` prepare parity with quiet, full-discover pure nested-tree parents, and planned headers under verbose

```
# Fix 1 — verbose prepare must not abort after Light succeeded
parent tree + intermediate/ (no SETUP) + nested DOCTEST
  quiet:  doctest test parent/        -> Light→Hydrate -> 1 case -> exit 0
  -v:     doctest test -v parent/     -> same 1 case -> exit 0
          -/-> intermediate/SETUP.md: must have a Go code block

# Fix 2 — full discover allows pure nested-tree parent dirs
DiscoverTreeCases(parent) when intermediate has no SETUP
  -> OK (1 parent leaf); nested DOCTEST skipped at boundary
DiscoverTreeCases when intermediate SETUP.md exists without Go
  -> still error: must have a Go code block

# Fix 3 — planned count always under -v
doctest test -v <tree>      -> stderr has (N tests) / planned before go test
doctest test -v ./...       -> doctest: workspace (N trees, M tests) then cd && go test
```

## Preconditions

- Nested root: module root is three levels above `d.DOCTEST_ROOT`
  (`tests/test/verbose-discover-parity` → workspace root).
- Shared doctest binary via `testbin.Ensure` for CLI leaves.
- Full-discover leaves call `core.DiscoverTreeCases` (no nested compile).
- Classic TDD: verbose prepare + pure-nested full discover + planned-header
  leaves are **RED** until implementer; quiet + SETUP-exists-without-Go stay **GREEN**.
- Child subprocess env only (`cmd.Env`); no parent `os.Setenv`.

## Steps

1. Resolve module root from `d.DOCTEST_ROOT`; build/reuse `req.Bin`.
2. Provide fixture writers: mega-style parent+nested intermediate, 1-pass tree,
   multi-tree workspace module, intermediate SETUP-without-Go.
3. Provide helpers to detect planned counts and forbidden intermediate errors.

## Context

- Quiet path already uses Light → filter → Hydrate (does not walk intermediate
  SETUP for non-ancestor dirs of selected leaves).
- Bug: under `-v`, prepare re-runs full `DiscoverTreeCasesVerbose` as a hard
  fail after Light succeeded, poisoning the parent tree when intermediate only
  holds a nested DOCTEST.
- Full discover currently treats every intermediate dir as requiring SETUP Go
  even when SETUP.md is absent.
- Workspace verbose path prints `cd … && go test` without the quiet planned line.

```go
import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

var (
	bt = string([]byte{96, 96, 96}) // triple backtick for nested markdown fixtures

	// Quiet single-tree: "doctest: <path> (N tests)"
	// Also accept "─── N test cases" and "N tests planned" variants under -v.
	plannedTestsRe = regexp.MustCompile(
		`(?i)(?:\((\d+)\s+tests?\)|───\s*(\d+)\s+test cases?|(\d+)\s+tests?\s+planned)`,
	)
	// Workspace: "doctest: workspace (N trees, M tests)" or hub label.
	plannedWorkspaceRe = regexp.MustCompile(
		`(?i)doctest:\s+workspace(?:\s+hub)?\s+\((\d+)\s+trees?,\s*(\d+)\s+tests?\)`,
	)
	goBlockErrRe = regexp.MustCompile(`(?i)must have a Go code block`)
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 180 * time.Second
	}
	if req.Op == "" || req.Op == "cli" || req.Op == "dual_cli" {
		req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", "..", ".."))
	}
	return nil
}

func doctestGoBlock(code string) string {
	return "## Test\n\n" + bt + "go\n" + code + "\n" + bt + "\n"
}

func scenarioFence(line string) string {
	return bt + "\n" + line + "\n" + bt
}

func tinyRunGo() string {
	return `import "testing"

type Request struct{}
type Response struct{}

func Run(t *testing.T, req *Request) (*Response, error) {
	_ = req
	return &Response{}, nil
}
`
}

func tinyRootSetupMD() string {
	return "# Scenario\n\n**Feature**: fixture root\n\n" + scenarioFence("root setup") +
		"\n\n## Steps\n1. no-op\n\n" +
		doctestGoBlock("import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n")
}

func tinyLeafSetupMD() string {
	return "# Scenario\n\n**Feature**: fixture leaf\n\n" + scenarioFence("leaf setup") +
		"\n\n## Steps\n1. no-op\n\n" +
		doctestGoBlock("import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n")
}

func tinyLeafAssertMD() string {
	return "# Scenario\n\n**Feature**: leaf passes\n\n" + scenarioFence("leaf pass") +
		"\n\n## Expected\n- pass\n\n" +
		doctestGoBlock("import \"testing\"\n\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n")
}

func writeTinyTree(t *testing.T, dir, leafName string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(tinyRunGo())), 0644); err != nil {
		t.Fatalf("write DOCTEST.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(tinyRootSetupMD()), 0644); err != nil {
		t.Fatalf("write root SETUP: %v", err)
	}
	leaf := filepath.Join(dir, leafName)
	if err := os.MkdirAll(leaf, 0755); err != nil {
		t.Fatalf("mkdir leaf: %v", err)
	}
	if err := os.WriteFile(filepath.Join(leaf, "SETUP.md"), []byte(tinyLeafSetupMD()), 0644); err != nil {
		t.Fatalf("write leaf SETUP: %v", err)
	}
	if err := os.WriteFile(filepath.Join(leaf, "ASSERT.md"), []byte(tinyLeafAssertMD()), 0644); err != nil {
		t.Fatalf("write leaf ASSERT: %v", err)
	}
}

// createMegaParentNestedFixture builds:
//
//	parent/
//	  DOCTEST.md + SETUP.md
//	  own_leaf/          ← 1 parent-tree case
//	  intermediate/      ← NO SETUP.md, NO ASSERT (pure nested-tree parent dir)
//	    nested/
//	      DOCTEST.md + SETUP.md + nested_leaf/
//
// Running doctest test on parent must see only own_leaf. Nested is a separate root.
// intermediate/ is the layout that full discover / -v re-walk wrongly rejects today.
func createMegaParentNestedFixture(t *testing.T) (moduleDir, parentDir string) {
	t.Helper()
	moduleDir = t.TempDir()
	if err := os.WriteFile(filepath.Join(moduleDir, "go.mod"), []byte("module example.com/vdiscparity\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	parentDir = filepath.Join(moduleDir, "parent")
	writeTinyTree(t, parentDir, "own_leaf")

	// Pure nested intermediate: no SETUP on intermediate itself.
	intermediate := filepath.Join(parentDir, "intermediate")
	if err := os.MkdirAll(intermediate, 0755); err != nil {
		t.Fatalf("mkdir intermediate: %v", err)
	}
	nested := filepath.Join(intermediate, "nested")
	writeTinyTree(t, nested, "nested_leaf")
	return moduleDir, parentDir
}

// createIntermediateSetupNoGoFixture is like mega parent, but intermediate has a
// prose-only SETUP.md (no Go block) — full discover must still error.
func createIntermediateSetupNoGoFixture(t *testing.T) (moduleDir, parentDir string) {
	t.Helper()
	moduleDir, parentDir = createMegaParentNestedFixture(t)
	intermediateSetup := filepath.Join(parentDir, "intermediate", "SETUP.md")
	body := "# Scenario\n\n**Feature**: bad intermediate SETUP without Go\n\n" +
		scenarioFence("intermediate setup missing go") +
		"\n\n## Steps\n1. intentionally no Go block\n"
	if err := os.WriteFile(intermediateSetup, []byte(body), 0644); err != nil {
		t.Fatalf("write intermediate SETUP: %v", err)
	}
	return moduleDir, parentDir
}

// createFastPassTree builds a 1-pass fixture for planned-header single-tree.
func createFastPassTree(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, tmp, []testtree.LeafSpec{
		{Name: "pass_0", Steps: "pass", Expected: "passes"},
	})
	return tmp
}

// createWorkspaceTwoTrees builds a module with two 1-leaf trees for ./... workspace.
func createWorkspaceTwoTrees(t *testing.T) string {
	t.Helper()
	mod := t.TempDir()
	if err := os.WriteFile(filepath.Join(mod, "go.mod"), []byte("module example.com/vdiscws\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	testtree.WritePassFailTree(t, filepath.Join(mod, "tree_a"), 1, 0)
	testtree.WritePassFailTree(t, filepath.Join(mod, "tree_b"), 1, 0)
	return mod
}

func combinedOutput(resp *Response) string {
	if resp == nil {
		return ""
	}
	return resp.Stdout + "\n" + resp.Stderr
}

// parsePlannedTests returns the first planned test count found in combined output, or -1.
func parsePlannedTests(out string) int {
	m := plannedTestsRe.FindStringSubmatch(out)
	if m == nil {
		return -1
	}
	for i := 1; i < len(m); i++ {
		if m[i] == "" {
			continue
		}
		n, err := strconv.Atoi(m[i])
		if err == nil {
			return n
		}
	}
	return -1
}

// parseWorkspacePlanned returns trees, tests from workspace planned line, or -1,-1.
func parseWorkspacePlanned(out string) (trees, tests int) {
	m := plannedWorkspaceRe.FindStringSubmatch(out)
	if m == nil {
		return -1, -1
	}
	tr, err1 := strconv.Atoi(m[1])
	te, err2 := strconv.Atoi(m[2])
	if err1 != nil || err2 != nil {
		return -1, -1
	}
	return tr, te
}

func hasIntermediateGoBlockError(out string) bool {
	if !goBlockErrRe.MatchString(out) {
		return false
	}
	// Path-sensitive: intermediate SETUP (with or without intermediate/ prefix).
	low := strings.ToLower(out)
	return strings.Contains(low, "intermediate") && strings.Contains(low, "setup.md")
}

func requireExit0(t *testing.T, resp *Response, err error, label string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: harness run failed: %v\nstdout:\n%s\nstderr:\n%s",
			label, err, resp.Stdout, resp.Stderr)
	}
	if resp == nil {
		t.Fatalf("%s: nil response", label)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("%s: expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s",
			label, resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}

```
