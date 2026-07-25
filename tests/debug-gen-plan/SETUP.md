# Scenario

**Feature**: DOCTEST_DEBUG gen-plan=1 prints generate plan/result trees on stderr

```
# parse
DOCTEST_DEBUG string
  -> debug.Parse (fail-closed)
  -> Settings.GenPlan / BypassGoTest

# product generate (isolated gen root)
DOCTEST_DEBUG=gen-plan=1,bypass-go-test=1
  + doctest test <fixture> --gen-dir <tmp>
  -> plan trees on stderr (arg[i/n], merged)
  -> generate (+ prune)
  -> result trees on stderr (status colors, summary)
```

## Preconditions

- Nested root under `tests/debug-gen-plan/` (firewall from parent `tests/`).
- Module root: `filepath.Join(d.DOCTEST_ROOT, "..", "..")`.
- **Classic TDD**: `gen-plan` not implemented — expect RED on acceptance + plan/result.
- Product leaves use small temp fixtures (`testtree.WriteMinimalRunnableTree`) and
  isolated `--gen-dir` / `DOCTEST_CACHE_HOME`.
- Prefer `cli.RunWithWriters` when no env injection; when `DebugEnv` is set use
  product binary + `cmd.Env` (Parallel-safe; never process Setenv).
- Session paths via `d *session.Doctest` only — never `os.Getenv("DOCTEST_SESSION_ID")`.

## Steps

1. Root Setup resolves module root and ensures product binary via `testbin.Ensure`.
2. Leaves set `Mode`, `DebugEnv`, fixtures, and CLI `Args`.
3. `Run` either parses DebugEnv or invokes CLI (once or twice for warm).
4. Assert inspects parse fields or stderr plan/result markers.

## Context

- Default DebugEnv for generate leaves: `gen-plan=1,bypass-go-test=1` (fast plan-only).
- Bookkeeping names under gen root: `go.mod`, `go.sum`, `doctest.gen-manifest`,
  `doctest.tidy-done`.
- Markers to match on stderr (flexible whitespace):
  - `doctest: DOCTEST_DEBUG gen-plan=1` (or equivalent banner)
  - `gen-plan: invocation`
  - `gen-plan: arg[i/n]`
  - `gen-plan: merged` (multi only)
  - `gen-plan: result`
  - `summary: modified=` / `unchanged=` / `deleted=`

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

const (
	debugGenPlanBypass = "gen-plan=1,bypass-go-test=1"
	genManifestName    = "doctest.gen-manifest"
	tidyDoneName       = "doctest.tidy-done"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	if req.Timeout <= 0 {
		req.Timeout = 120 * time.Second
	}
	// Product binary for Env-injected leaves; parse leaves ignore Bin.
	modRoot := moduleRootFrom(d)
	req.Bin = testbin.Ensure(t, modRoot)
	return nil
}

func moduleRootFrom(d *session.Doctest) string {
	return filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
}

// isolateEnv returns child-only env extras (no DOCTEST_DEBUG — Run adds that).
func isolateEnv(t *testing.T) []string {
	t.Helper()
	cacheHome := t.TempDir()
	gocache := t.TempDir()
	return []string{
		"GOWORK=off",
		"DOCTEST_CACHE_HOME=" + cacheHome,
		"GOCACHE=" + gocache,
	}
}

// prepareSingleTree writes one minimal runnable doctest tree in a temp module.
// Returns fixture root (DOCTEST.md parent). WorkDir is the module root.
func prepareSingleTree(t *testing.T, req *Request) {
	t.Helper()
	work := t.TempDir()
	if err := os.WriteFile(filepath.Join(work, "go.mod"), []byte("module example.com/gen-plan-single\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	tree := filepath.Join(work, "tree")
	if err := os.MkdirAll(tree, 0o755); err != nil {
		t.Fatal(err)
	}
	testtree.WriteMinimalRunnableTree(t, tree, []testtree.LeafSpec{
		{Name: "leaf", Steps: "pass", Expected: "ok"},
	})
	req.WorkDir = work
	req.ModuleRoot = work
	req.FixtureDir = tree
	req.TreeRoot = tree
	req.GenDir = filepath.Join(t.TempDir(), "gen")
	if err := os.MkdirAll(req.GenDir, 0o755); err != nil {
		t.Fatal(err)
	}
	req.Env = isolateEnv(t)
}

// prepareMultiArgTwoTrees writes module with tree-a + tree-b for multi-arg CLI.
func prepareMultiArgTwoTrees(t *testing.T, req *Request) {
	t.Helper()
	work := t.TempDir()
	if err := os.WriteFile(filepath.Join(work, "go.mod"), []byte("module example.com/gen-plan-multi\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	treeA := filepath.Join(work, "tree-a")
	treeB := filepath.Join(work, "tree-b")
	if err := os.MkdirAll(treeA, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(treeB, 0o755); err != nil {
		t.Fatal(err)
	}
	testtree.WriteMinimalRunnableTree(t, treeA, []testtree.LeafSpec{
		{Name: "a1", Steps: "a1", Expected: "ok"},
	})
	testtree.WriteMinimalRunnableTree(t, treeB, []testtree.LeafSpec{
		{Name: "b1", Steps: "b1", Expected: "ok"},
	})
	req.WorkDir = work
	req.ModuleRoot = work
	req.TreeRoot = treeA
	req.TreeRootB = treeB
	req.FixtureDir = treeA
	req.GenDir = filepath.Join(t.TempDir(), "gen")
	if err := os.MkdirAll(req.GenDir, 0o755); err != nil {
		t.Fatal(err)
	}
	req.Env = isolateEnv(t)
}

// baseTestArgs builds common args: test [color flags] --gen-dir GenDir -count=1 paths…
func baseTestArgs(req *Request, colorFlag string, paths ...string) []string {
	args := []string{"test"}
	if colorFlag != "" {
		args = append(args, colorFlag)
	}
	args = append(args, "--gen-dir", req.GenDir, "-count=1")
	args = append(args, paths...)
	return args
}

func stderrHas(s, sub string) bool {
	return strings.Contains(s, sub)
}

// extractGenPlanSection returns lines from the first marker through end or next major section.
func extractGenPlanSection(stderr, marker string) string {
	idx := strings.Index(stderr, marker)
	if idx < 0 {
		return ""
	}
	return stderr[idx:]
}

func countSubstr(s, sub string) int {
	if sub == "" {
		return 0
	}
	return strings.Count(s, sub)
}
```
