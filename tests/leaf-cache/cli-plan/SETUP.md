# Scenario

**Feature**: multi-arg `treeA treeB` leaf-cache product path (**L3 e2e**, `label: e2e`)

```
# P3 product path (CLI one plan)
module with tree-a/ + tree-b/ (two DOCTEST.md roots)
  -> doctest test tree-a tree-b     # multi-arg union of roots
  -> prepare → bind → run → report
  -> leaf-cache PutPass / warm GetPass skip / summary Cached
  # same policy as:
  #   doctest test tree-a           (single)
  #   doctest test <mod>/...        (workspace)
```

## Preconditions

- **Layer L3** — nested multi-arg product; every leaf has `label: e2e`.
- Reuses the same multi-tree module fixture as `workspace/` (`prepareWorkspaceAllPass`).
- Nested CLI: `Op=runtime_multi`, isolated `DOCTEST_CACHE_HOME` / `DOCTEST_LEAF_CACHE`,
  fresh `GOCACHE` per invocation (parent helpers).
- Args use **two explicit roots** (`test tree-a tree-b`), **not** `<mod>/...`
  (that path is sealed under `workspace/`).
- Skipped by default discovery; run with `--label e2e`.
- **Shared-policy product (GREEN)**: multi-arg may fan to N× `TestWithStats` or
  a shared plan builder; asserts check **observable** Cached / `-count` policy
  parity with single-tree and workspace, not call-graph shape.

## Steps

1. Resolve module root; ensure binary; set `Op=runtime_multi`, timeout, isolate env.
2. Child leaves write the multi-tree fixture and multi-arg Args/Args2/Args3.
3. Assert summary Cached / exit codes for warm and count bypass.

## Context

- **Significance**: third product invocation shape — multi-arg CLI. Single-tree
  and workspace `./...` leaf-cache are sealed under `runtime/**` and
  `workspace/**`. Together they prove leaf-cache is multi-path, not single-tree only.
- Do not re-test library PreparePassPlan (that's `runsuite/**`).
- Summary **Cached** is leaf-cache product skips (fresh GOCACHE in harness).

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

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	req.Timeout = 240 * time.Second
	// tests/leaf-cache -> tests -> module root
	modRoot := filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	req.Bin = testbin.Ensure(t, modRoot)
	req.Env = isolateRuntimeEnv(t)
	return nil
}

// multiArgTwoTrees returns Args for explicit multi-root invocation:
//   doctest test <tree-a> <tree-b>
// TreeRoot / TreeRootB must already be set (e.g. prepareMultiArgAllPass).
func multiArgTwoTrees(req *Request) []string {
	return []string{"test", req.TreeRoot, req.TreeRootB}
}

// prepareMultiArgTwoTrees writes a single-module multi-tree fixture for multi-arg:
//
//	$WorkDir/
//	  go.mod
//	  tree-a/   DOCTEST + leavesA
//	  tree-b/   DOCTEST + leavesB
//
// Same layout as workspace/ prepareWorkspaceTwoTrees (sibling branch helpers
// are not inherited — keep a local copy for the cli-plan path).
func prepareMultiArgTwoTrees(t *testing.T, req *Request, leavesA, leavesB []testtree.LeafSpec) {
	t.Helper()
	work := t.TempDir()
	if err := os.WriteFile(filepath.Join(work, "go.mod"), []byte("module example.com/cli-plan-leafcache\n\ngo 1.22\n"), 0o644); err != nil {
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
	if len(leavesA) == 0 {
		leavesA = []testtree.LeafSpec{{Name: "a1", Steps: "a1", Expected: "ok"}}
	}
	if len(leavesB) == 0 {
		leavesB = []testtree.LeafSpec{{Name: "b1", Steps: "b1", Expected: "ok"}}
	}
	testtree.WriteMinimalRunnableTree(t, treeA, leavesA)
	testtree.WriteMinimalRunnableTree(t, treeB, leavesB)
	req.WorkDir = work
	req.FixtureDir = work
	req.ModuleRoot = work
	req.TreeRoot = treeA
	req.TreeRootB = treeB
	req.FixtureB = treeB
	req.LeafDir = filepath.Join(treeA, leavesA[0].Name)
	req.LeafDirB = filepath.Join(treeB, leavesB[0].Name)
}

// prepareMultiArgAllPass is two trees, one pass leaf each (warm happy path).
func prepareMultiArgAllPass(t *testing.T, req *Request) {
	t.Helper()
	prepareMultiArgTwoTrees(t, req,
		[]testtree.LeafSpec{{Name: "a1", Steps: "a1", Expected: "ok"}},
		[]testtree.LeafSpec{{Name: "b1", Steps: "b1", Expected: "ok"}},
	)
}
```
