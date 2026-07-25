# Scenario

**Feature**: gen manifest / gen-plan scope and lifecycle policies

Covers subset isolation, warm stability, unmanaged plants, managed missing
rewrite, out-of-scope non-heal, source change → modified, source leaf remove →
deleted, `-a` cold rewrite, and go.mod bookkeeping recreate.

## Preconditions

- Nested root `tests/gen-manifest-scope/`.
- Isolated `--gen-dir` + cache env; product binary via `testbin.Ensure`.

## Steps

1. Root Setup ensures binary.
2. Leaves prepare fixtures and phase hooks (`DeleteGenRels`, `PlantRel`, …).
3. Run executes up to three CLI phases; Assert checks disk, ledger, gen-plan.

```go
import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

const debugGenPlanBypass = "gen-plan=1,bypass-go-test=1"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	if req.Timeout <= 0 {
		req.Timeout = 120 * time.Second
	}
	modRoot := filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	req.Bin = testbin.Ensure(t, modRoot)
	return nil
}

func isolateEnv(t *testing.T) []string {
	t.Helper()
	return []string{
		"GOWORK=off",
		"DOCTEST_CACHE_HOME=" + t.TempDir(),
		"GOCACHE=" + t.TempDir(),
	}
}

func prepareTwoTreeModule(t *testing.T, req *Request) {
	t.Helper()
	work := t.TempDir()
	if err := os.WriteFile(filepath.Join(work, "go.mod"), []byte("module example.com/gen-manifest-scope\n\ngo 1.22\n"), 0o644); err != nil {
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
	req.TreeA = treeA
	req.TreeB = treeB
	req.GenDir = filepath.Join(t.TempDir(), "gen")
	if err := os.MkdirAll(req.GenDir, 0o755); err != nil {
		t.Fatal(err)
	}
	req.Env = isolateEnv(t)
}

// prepareSingleTreeModule writes go.mod + tree/ with one leaf named "leaf".
func prepareSingleTreeModule(t *testing.T, req *Request) {
	t.Helper()
	prepareSingleTreeModuleWithLeaves(t, req, []testtree.LeafSpec{
		{Name: "leaf", Steps: "pass", Expected: "ok"},
	})
}

// prepareSingleTreeModuleWithLeaves allows multi-leaf single trees (orphan delete).
func prepareSingleTreeModuleWithLeaves(t *testing.T, req *Request, leaves []testtree.LeafSpec) {
	t.Helper()
	work := t.TempDir()
	if err := os.WriteFile(filepath.Join(work, "go.mod"), []byte("module example.com/gen-manifest-scope-single\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	tree := filepath.Join(work, "tree")
	if err := os.MkdirAll(tree, 0o755); err != nil {
		t.Fatal(err)
	}
	testtree.WriteMinimalRunnableTree(t, tree, leaves)
	req.WorkDir = work
	req.TreeA = tree
	req.GenDir = filepath.Join(t.TempDir(), "gen")
	if err := os.MkdirAll(req.GenDir, 0o755); err != nil {
		t.Fatal(err)
	}
	req.Env = isolateEnv(t)
}

func baseArgs(req *Request, paths ...string) []string {
	args := []string{"test", "--gen-dir", req.GenDir, "-count=1", "--no-color"}
	return append(args, paths...)
}

// baseArgsForceA is like baseArgs but inserts -a after "test".
func baseArgsForceA(req *Request, paths ...string) []string {
	args := []string{"test", "-a", "--gen-dir", req.GenDir, "-count=1", "--no-color"}
	return append(args, paths...)
}
```
