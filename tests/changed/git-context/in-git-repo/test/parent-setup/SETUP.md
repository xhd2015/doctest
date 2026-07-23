# Scenario

**Feature**: changing a parent `SETUP.md` selects all descendant leaves

```
# shared/SETUP.md in changed list
FilterByChangedFiles -> [shared/leaf_a, shared/leaf_b]
```

## Steps

1. Create shared-parent two-leaf tree.
2. Set changed path to `shared/SETUP.md`.
3. Run filter policy.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	fx := createSharedParentTwoLeafTree(t)
	applyPolicyBase(req, fx)
	req.ChangedFiles = []string{treeRel(fx, "shared", "SETUP.md")}
	return nil
}
```
