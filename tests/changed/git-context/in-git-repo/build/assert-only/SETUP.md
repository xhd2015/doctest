# Scenario

**Feature**: build selection for a single ASSERT change matches test filter policy

```
# same FilterByChangedFiles as test assert-only
changed leaf_a/ASSERT.md -> [leaf_a]
```

## Steps

1. Create flat two-leaf tree.
2. Set changed ASSERT path.
3. Assert filtered set (no compile required for selection policy).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	applyPolicyBase(req, fx)
	req.ChangedFiles = []string{treeRel(fx, "leaf_a", "ASSERT.md")}
	return nil
}
```
