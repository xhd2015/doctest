# Scenario

**Feature**: changing one leaf `ASSERT.md` selects only that leaf

```
# only leaf_a ASSERT.md in changed list
changed leaf_a/ASSERT.md -> FilterByChangedFiles -> [leaf_a]
```

## Steps

1. Create flat two-leaf tree.
2. Set `ChangedFiles` to `leaf_a/ASSERT.md` under the tree.
3. Run filter policy.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	applyPolicyBase(req, fx)
	req.ChangedFiles = []string{treeRel(fx, "leaf_a", "ASSERT.md")}
	return nil
}
```
