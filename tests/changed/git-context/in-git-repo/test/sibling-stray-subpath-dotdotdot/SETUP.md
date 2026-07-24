# Scenario

**Feature**: subpath discovery still ignores sibling stray files

```
# same filter policy as sibling-stray-untracked
changed leaf_a/ASSERT.md + leaf_b/stray.go -> [leaf_a]
```

## Steps

1. Create flat two-leaf tree.
2. Set changed list including sibling stray.
3. Assert only leaf_a.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	applyPolicyBase(req, fx)
	req.ChangedFiles = []string{
		treeRel(fx, "leaf_a", "ASSERT.md"),
		treeRel(fx, "leaf_b", "stray.go"),
	}
	return nil
}
```
