# Scenario

**Feature**: mixed tree runs fast leaf and skips labeled leaf

```
discovery FilterCasesByLabel → run {fast_leaf}; skip {labeled_leaf}
```

## Steps

1. Create mixed temp tree.
2. Default discovery options (no LabelExprs, not ExplicitLeaf).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabeledTree(t, true, "ui-automation", "heavy ui test")
	return nil
}
```
