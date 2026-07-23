# Scenario

**Feature**: a new leaf's files select only that leaf

```
# new leaf_c ASSERT + SETUP in changed list; baseline leaves unchanged
FilterByChangedFiles -> [leaf_c]
```

## Steps

1. Create flat two-leaf tree and add `leaf_c` on disk (so discovery finds it).
2. Set changed paths to leaf_c markdown only.
3. Assert only leaf_c is selected.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	writeLeaf(t, fx.TreeDir, "leaf_c")
	applyPolicyBase(req, fx)
	req.ChangedFiles = []string{
		treeRel(fx, "leaf_c", "ASSERT.md"),
		treeRel(fx, "leaf_c", "SETUP.md"),
	}
	return nil
}
```
