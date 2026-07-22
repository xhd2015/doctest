# Scenario

**Feature**: edit shared package used by two leaves → only leaf-d key stays stable

```
keys0 = ComputeLeafKey(leaf-ab-1, leaf-ab-2, leaf-d)
edit shared/a Version
keys1 = ComputeLeafKey(...)
# leaf-ab-* keys change; leaf-d key unchanged
```

## Preconditions

- Same partial-package fixture as sibling leaf.
- Mutation `polish_edit_shared_a`.

## Steps

1. Build fixture.
2. Op=`partial_package_keys` with shared/a mutation.

## Context

- Library equivalent of "1 Cached" product warm: alone leaf keeps prior key.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "partial_package_keys"
	_ = preparePartialPackageDepsFixture(t, req)
	req.Mutation = "polish_edit_shared_a"
	return nil
}
```
