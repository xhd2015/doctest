# Scenario

**Feature**: changing root `DOCTEST.md` selects all leaves in the tree

```
# root DOCTEST.md in changed list
FilterByChangedFiles -> all leaves
```

## Steps

1. Create flat two-leaf tree.
2. Set changed path to root `DOCTEST.md`.
3. Run filter policy.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	applyPolicyBase(req, fx)
	req.ChangedFiles = []string{treeRel(fx, "DOCTEST.md")}
	return nil
}
```
