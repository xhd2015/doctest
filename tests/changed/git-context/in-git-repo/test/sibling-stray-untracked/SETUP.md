# Scenario

**Feature**: unrelated non-doctest files in a sibling leaf must not widen selection

```
# leaf_a ASSERT + leaf_b/stray.go both "changed"
FilterByChangedFiles -> [leaf_a] only
```

## Steps

1. Create flat two-leaf tree (stray path need not exist on disk for path math).
2. Set changed list to ASSERT + sibling stray.go.
3. Assert only leaf_a is selected.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	applyPolicyBase(req, fx)
	req.ChangedFiles = []string{
		treeRel(fx, "leaf_a", "ASSERT.md"),
		treeRel(fx, "leaf_b", "stray.go"),
	}
	return nil
}
```
