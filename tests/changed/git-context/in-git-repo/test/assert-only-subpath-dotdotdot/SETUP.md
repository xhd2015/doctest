# Scenario

**Feature**: `./tests/...` discovery uses the same filter once a tree root is known

```
# CLI argv ./tests/... is orthogonal to selection policy
changed leaf_a/ASSERT.md -> FilterByChangedFiles -> [leaf_a]
```

## Steps

1. Create flat two-leaf tree under `tests/fixture` (matches subpath layout).
2. Set synthetic changed ASSERT path.
3. Assert identical selection to assert-only.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	applyPolicyBase(req, fx)
	req.ChangedFiles = []string{treeRel(fx, "leaf_a", "ASSERT.md")}
	return nil
}
```
