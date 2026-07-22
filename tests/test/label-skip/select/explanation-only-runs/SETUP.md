# Scenario

**Feature**: explanation-only frontmatter does not trigger skip

```
explanation without label → leaf runs under discovery
```

## Steps

1. Create explanation-only tree.

```go
func Setup(t *testing.T, req *Request) error {
	req.TreeRoot = writeExplanationOnlyTree(t, "documentation note only")
	return nil
}
```
