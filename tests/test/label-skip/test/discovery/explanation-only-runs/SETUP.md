# Scenario

**Feature**: explanation-only frontmatter does not trigger skip

```
# explanation without label
doctest test <tree-root> -> runs leaf, no SKIPPED summary
```

## Steps

1. Create temp tree with explanation-only frontmatter.
2. Run `doctest test <tree-root>`.

```go
func Setup(t *testing.T, req *Request) error {
	root := writeExplanationOnlyTree(t, "documentation note only")
	req.Args = []string{"test", root}
	return nil
}
```