# Scenario

**Feature**: partial leaf-cache key DAG across 3 leaves / 4 local packages (**L2 in-process**)

```
# packages
shared/{a,b,c} -> leaf-ab-1, leaf-ab-2
alone/d        -> leaf-d only

# library (no product binary)
ComputeLeafKey per leaf before/after package edit
  -> edit alone/d invalidates leaf-d only (shared leaves stable)
  -> edit shared/a invalidates leaf-ab-* only (leaf-d stable)
```

## Preconditions

- **Layer L2** — `ComputeLeafKey` only; no `testbin` / nested `doctest test`.
- Fixture: 4 packages + 3 leaves under a temp module (`preparePartialPackageDepsFixture`).
- Unlabeled (fast discovery).

## Steps

1. Children call `preparePartialPackageDepsFixture` and set `Mutation`.
2. Run `Op=partial_package_keys` (keys before → mutate → keys after).
3. Assert which leaf-key groups stayed stable (`Hit` / `HitB`).

## Context

- Proves per-leaf import-closure hashing, not whole-tree wipe.
- Product warm **Cached** for multi-leaf remains under L3 `runtime/**` / `polish/**`
  only where still e2e; this branch is the library mass for partial deps.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	if req.GoVersion == "" {
		req.GoVersion = "go1.25.0"
	}
	req.Op = "partial_package_keys"
	return nil
}
```
