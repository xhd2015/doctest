# Scenario

**Feature**: editing a local package imported by the spine changes the leaf key

```
key0 = ComputeLeafKey(leaf importing helper)
mutate pkg/helper/helper.go
key1 = ComputeLeafKey(...)
# key0 != key1
```

## Preconditions

- Fixture module with `pkg/helper` imported by leaf ASSERT.
- Mutation `local_imported` (same effect as polish_edit_local_dep).

## Steps

1. `prepareLocalDepPassFixture`.
2. Op=`compute_mutate` with `local_imported`.

## Context

- Library invalidation; product "0 Cached after dep edit" is the L3 consequence of a new key.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "compute_mutate"
	_ = prepareLocalDepPassFixture(t, req)
	req.Mutation = "local_imported"
	return nil
}
```
