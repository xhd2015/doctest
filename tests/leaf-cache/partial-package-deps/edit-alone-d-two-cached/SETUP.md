# Scenario

**Feature**: edit package only leaf-d uses → peer leaf keys stay stable

```
keys0 = ComputeLeafKey(leaf-ab-1, leaf-ab-2, leaf-d)
edit alone/d Version
keys1 = ComputeLeafKey(...)
# leaf-ab-1 and leaf-ab-2 keys unchanged; leaf-d key changes
```

## Preconditions

- `preparePartialPackageDepsFixture` layout.
- Mutation `polish_edit_alone_d`.

## Steps

1. Build fixture.
2. Op=`partial_package_keys` with alone/d mutation.

## Context

- Library equivalent of "2 Cached" product warm: shared leaves keep prior keys.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "partial_package_keys"
	_ = preparePartialPackageDepsFixture(t, req)
	req.Mutation = "polish_edit_alone_d"
	return nil
}
```
