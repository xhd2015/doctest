# Scenario

**Feature**: unlabeled-only tree runs without skips

```
FilterCasesByLabel → run {plain_leaf}; skip {}
```

## Steps

1. Create unlabeled tree.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeUnlabeledTree(t)
	return nil
}
```
