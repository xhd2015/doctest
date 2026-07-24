# Scenario

**Feature**: multi-prep PassPlan — tree-qualified skip list + partial record

```
# prepare (skip enabled, both warm)
PreparePassPlan([A/leaf, B/leaf], skipEnabled=true)
  -> Keys has 2 tree-qualified identities
  -> Skip contains both identities (sorted)

# prepare (only A warm)
  -> Skip contains only identity(A)

# prepare (skip disabled)
  -> Skip empty even if store warm

# record (partial fail)
RecordPasses(keys, failed={idA}, allPassed=false)
  -> GetPass(keyA)=false, GetPass(keyB)=true
```

## Preconditions

- Twin trees + isolated StoreRoot.
- Leaves seed store via PutPass of ComputeLeafKey before prepare when testing warm paths.
- `req.PrepWarm` controls which trees are pre-warmed: `both` | `a` | `none`.
- `req.SkipEnabled` defaults true; false models `-count` / `-a` / `--no-leaf-cache`.

## Steps

1. prepareTwinTrees; open store under StoreRoot.
2. Optionally PutPass keys for warm trees.
3. Leaf sets Op + PrepWarm + SkipEnabled / FailedSide.

## Context

- These helpers are the extract seam single-tree `prepareLeafCache` /
  `recordLeafCachePasses` should share so multi-tree RunSuite can call one plan
  without identity collisions.
- Not an e2e workspace `./...` test (P2).

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/leafcache"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareTwinTrees(t, req)
	if req.StoreRoot == "" {
		req.StoreRoot = filepath.Join(t.TempDir(), "leaf-cache-v1")
	}
	if req.GoVersion == "" {
		req.GoVersion = "go1.25.0"
	}
	// Default: skip enabled (product default when no -count/-a/--no-leaf-cache).
	if !req.SkipEnabledSet {
		req.SkipEnabled = true
		req.SkipEnabledSet = true
	}
	return nil
}

// seedWarmStore PutPass's ComputeLeafKey for tree A and/or B per PrepWarm.
// Seeding only needs ComputeLeafKey + Store (already implemented).
func seedWarmStore(t *testing.T, req *Request) error {
	t.Helper()
	st, err := leafcache.NewStore(req.StoreRoot)
	if err != nil {
		return err
	}
	warmA := req.PrepWarm == "both" || req.PrepWarm == "a"
	warmB := req.PrepWarm == "both"
	if warmA {
		k, err := leafcache.ComputeLeafKey(leafcache.KeyInput{
			ModuleRoot: req.ModuleRoot,
			TreeRoot:   req.TreeRoot,
			LeafDir:    req.LeafDir,
			GoVersion:  req.GoVersion,
		})
		if err != nil {
			return err
		}
		if err := st.PutPass(k); err != nil {
			return err
		}
	}
	if warmB {
		k, err := leafcache.ComputeLeafKey(leafcache.KeyInput{
			ModuleRoot: req.ModuleRootB,
			TreeRoot:   req.TreeRootB,
			LeafDir:    req.LeafDirB,
			GoVersion:  req.GoVersion,
		})
		if err != nil {
			return err
		}
		if err := st.PutPass(k); err != nil {
			return err
		}
	}
	return nil
}
```
