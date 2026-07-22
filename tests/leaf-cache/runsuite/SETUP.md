# Scenario

**Feature**: extractable RunSuite leaf-cache helpers — multi-tree identity + prepare/record

```
# P1 extract seam (library; not full workspace ./... wiring)
twins treeA + treeB, same rel leaf path
  -> FormatLeafIdentity(tree, rel)  # tree-qualified tokens
  -> PreparePassPlan(store, leaves, goVer, skipEnabled)
       keys: map[identity]storeKey
       skip: []identity for GetPass hits when enabled
  -> RecordPasses(store, keys, failed, allPassed)
       PutPass only non-failed identities
```

## Preconditions

- Nested root: package `github.com/xhd2015/doctest/libdoc/leafcache`.
- Multi-prep helpers are **GREEN** / exported (`FormatLeafIdentity`,
  `PreparePassPlan`, `RecordPasses`).
- Twin fixtures via `prepareTwinTrees` (content-identical; different abs roots).
- Store roots always under `t.TempDir()` — never the user cache.
- End-to-end workspace / multi-arg product paths live under parent
  `workspace/**` and `cli-plan/**` (out of scope for this nested root).

## Steps

1. Root Setup defaults GoVersion and StoreRoot.
2. Group/leaf Setup chooses Op and warms the store when needed.
3. Run calls FormatLeafIdentity / PreparePassPlan / RecordPasses.
4. Assert tree-qualified uniqueness, skip membership, and partial PutPass.

## Context

### Expected public API (multi-prep extract — sealed)

```text
func FormatLeafIdentity(treeRoot, leafRel string) string
type LeafRef struct { TreeRoot, LeafRel string }
type PassPlan struct {
    Keys map[string]string
    Skip []string
}
func PreparePassPlan(store *Store, leaves []LeafRef, goVersion string, skipEnabled bool) (PassPlan, error)
func RecordPasses(store *Store, keys map[string]string, failed map[string]bool, allPassed bool)
```

### Why identity ≠ store key

- **Store key** already mixes abs TreeRoot into `ComputeLeafKey`.
- **Identity** is the in-memory / env token for skip lists and fail maps when
  multiple trees share one suite process. Bare relative paths collide.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	if req.GoVersion == "" {
		req.GoVersion = "go1.25.0"
	}
	if req.StoreRoot == "" {
		req.StoreRoot = filepath.Join(t.TempDir(), "leaf-cache-v1")
	}
	return nil
}

// prepareTwinTrees writes two content-identical single-leaf trees under work.
// Same relative path "leaf"; different absolute TreeRoots.
func prepareTwinTrees(t *testing.T, req *Request) {
	t.Helper()
	work := t.TempDir()
	writeOneLeafTree := func(root string) {
		testtree.WriteMinimalRunnableTree(t, root, []testtree.LeafSpec{
			{
				Name: "leaf",
				AssertGo: `import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	// twin_marker
	_ = req
	_ = resp
	_ = err
}`,
			},
		})
	}
	treeA := filepath.Join(work, "treeA")
	treeB := filepath.Join(work, "treeB")
	if err := os.MkdirAll(treeA, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(treeB, 0o755); err != nil {
		t.Fatal(err)
	}
	writeOneLeafTree(treeA)
	writeOneLeafTree(treeB)
	req.WorkDir = work
	req.TreeRoot = treeA
	req.LeafDir = filepath.Join(treeA, "leaf")
	req.ModuleRoot = treeA
	req.TreeRootB = treeB
	req.LeafDirB = filepath.Join(treeB, "leaf")
	req.ModuleRootB = treeB
}
```
