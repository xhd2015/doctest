# Scenario

**Feature**: all-labeled tree skips every leaf under discovery

```
FilterCasesByLabel → run {}; skip {labeled_leaf}
```

## Steps

1. Create labeled-only tree.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabeledTree(t, false, "human-guided-ui-test", "manual only")
	return nil
}
```
