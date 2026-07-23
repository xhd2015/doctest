# Scenario

**Feature**: only non-doctest file changes yield zero leaf selection

```
# README.md outside the tree is "changed"
FilterByChangedFiles -> [] (path not under doctest root)
```

## Steps

1. Create flat two-leaf tree.
2. Set changed path to repo-root `README.md`.
3. Run filter policy.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	applyPolicyBase(req, fx)
	req.ChangedFiles = []string{"README.md"}
	return nil
}
```
