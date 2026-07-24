# Scenario

**Feature**: LabelAll runs fast and labeled leaves

```
Options{LabelAll: true} → run both leaves
```

## Steps

1. Mixed tree; LabelAll = true.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabeledTree(t, true, "ui-automation", "heavy ui test")
	req.LabelAll = true
	return nil
}
```
