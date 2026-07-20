# Scenario

**Feature**: disk pass store GetPass / PutPass under an isolated root

```
# empty store
GetPass(key) -> false

# after explicit PutPass
PutPass(key) -> GetPass(key)=true

# separate Root
PutPass on A is invisible to store B
```

## Preconditions

- Store roots are always `t.TempDir()` (or subdirs thereof).
- Keys may be synthetic hex strings; ComputeLeafKey is optional for store-only leaves.
- PutPass is never implied by failure (no suite wiring in P1).

## Steps

1. Leaf sets Op to `store_put_get`, `store_missing`, or `store_isolate`.
2. Root Run opens `leafcache.NewStore(root)` and exercises Get/Put.
3. Assert Hit / HitB booleans.

## Context

- Product default `$CacheHome/doctest/leaf-cache/v1` is conceptual only; tests never write there.
- Significance: store is the durable side of P1; keys come from the key/ branch.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	if req.StoreRoot == "" {
		req.StoreRoot = filepath.Join(t.TempDir(), "leaf-cache-v1")
	}
	return nil
}
```
