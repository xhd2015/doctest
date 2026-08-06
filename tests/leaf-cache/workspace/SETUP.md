# Scenario

**Feature**: multi-tree `__workspace` leaf-cache product path (**L3 e2e**, `label: e2e`)

```
# P2 product path (not library-only multi-prep)
module with tree-a/ + tree-b/ (two DOCTEST.md roots)
  -> doctest test <mod>/...
  -> PrepareTree ×N + RunWorkspace / finishWorkspaceGoTest
  -> multi-prep PreparePassPlan + DOCTEST_LEAF_CACHE_SKIP_PATHS
  -> go test __workspace/suite
  -> RecordPasses + summary Cached = leaf-cache skip count
```

## Preconditions

- **Layer L3** — nested multi-tree product; every leaf has `label: e2e`.
- Product **GREEN**: `finishWorkspaceGoTest` prepares skip env, records passes,
  and reports summary **Cached** as programmatic leaf-cache skip count.
- Nested CLI: same `runtime_multi` Op as `runtime/**` (fresh GOCACHE per run).
- Fixture is a temp **module** with `go.mod` and **two** sibling DOCTEST roots
  so `test <mod>/...` takes the `testDotDotDotWorkspace` path (single-gen
  `__workspace`), not per-tree `TestWithStats`.
- Isolated `DOCTEST_CACHE_HOME` / `DOCTEST_LEAF_CACHE`; never the user cache.
- Requires built doctest binary (`testbin.Ensure`).
- Skipped by default discovery; run with `--label e2e`.

## Steps

1. Resolve module root; ensure binary; set `Op=runtime_multi`, timeout, isolate env.
2. Child leaves write the multi-tree fixture and Args/Args2/Args3.
3. Assert summary Cached / exit codes for warm, partial-fail, count bypass, isolation.

## Context

- **Significance**: product surface for multi-tree workspace — proves leaf-cache
  is not single-tree only (`doctest test <mod>/...` / `./...`).
- Nested `runsuite/**` seals FormatLeafIdentity / PreparePassPlan / RecordPasses
  as library helpers; this branch seals **end-to-end workspace wiring**.
- Hub multi-mod path is the same `finishWorkspaceGoTest` entry; single-gen
  workspace is the primary fixture (no separate hub leaf — same finish path).
- Skip tokens must be tree-qualified (or otherwise identity-safe) so same
  relative leaf paths across trees do not false-skip — see `isolation/`.
- Summary **Cached** here is leaf-cache product skips, not go testcache.

```go
import (
	"fmt"
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

// workspaceDotDotDot returns an absolute <mod>/... pattern for multi-tree workspace.
func workspaceDotDotDot(modRoot string) string {
	return filepath.Clean(modRoot) + "/..."
}

// prepareWorkspaceTwoTrees writes a single-module multi-tree fixture:
//
//	$WorkDir/
//	  go.mod
//	  tree-a/   DOCTEST + leavesA
//	  tree-b/   DOCTEST + leavesB
//
// FixtureDir/WorkDir is the module root; use workspaceDotDotDot for Args.
func prepareWorkspaceTwoTrees(t *testing.T, req *Request, leavesA, leavesB []testtree.LeafSpec) {
	t.Helper()
	work := t.TempDir()
	if err := os.WriteFile(filepath.Join(work, "go.mod"), []byte("module example.com/ws-leafcache\n\ngo 1.22\n"), 0o644); err != nil {
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

// prepareWorkspaceAllPass is two trees, one pass leaf each (warm happy path).
func prepareWorkspaceAllPass(t *testing.T, req *Request) {
	t.Helper()
	prepareWorkspaceTwoTrees(t, req,
		[]testtree.LeafSpec{{Name: "a1", Steps: "a1", Expected: "ok"}},
		[]testtree.LeafSpec{{Name: "b1", Steps: "b1", Expected: "ok"}},
	)
}

// prepareWorkspacePartialFail is tree-a pass + tree-b fail (multi-tree partial fail).
func prepareWorkspacePartialFail(t *testing.T, req *Request) {
	t.Helper()
	prepareWorkspaceTwoTrees(t, req,
		[]testtree.LeafSpec{{Name: "a_pass", Steps: "pass", Expected: "ok"}},
		[]testtree.LeafSpec{{
			Name:     "b_fail",
			Steps:    "fail",
			Expected: "fails",
			AssertGo: `import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Fatal("forced workspace fail leaf")
}`,
		}},
	)
}

// prepareWorkspaceSameRelpathPassFail writes twins with relative path "leaf":
// tree-a always passes; tree-b always fails. Used for cross-tree skip isolation.
func prepareWorkspaceSameRelpathPassFail(t *testing.T, req *Request) {
	t.Helper()
	prepareWorkspaceTwoTrees(t, req,
		[]testtree.LeafSpec{{
			Name: "leaf",
			AssertGo: `import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	// twin_ws_a
	_ = req
	_ = resp
	_ = err
}`,
		}},
		[]testtree.LeafSpec{{
			Name: "leaf",
			AssertGo: `import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Fatal("tree-b must execute body; false warm skip would hide this fail")
}`,
		}},
	)
}

// mustWorkspacePattern documents that Args must use /... for RunWorkspace.
func mustWorkspacePattern(modRoot string) string {
	p := workspaceDotDotDot(modRoot)
	if p == "" {
		panic(fmt.Sprintf("empty workspace pattern for %q", modRoot))
	}
	return p
}
```
