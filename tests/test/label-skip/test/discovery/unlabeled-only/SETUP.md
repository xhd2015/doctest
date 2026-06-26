# Scenario

**Feature**: unlabeled-only tree runs without skip summary

```
# no frontmatter labels anywhere
doctest test <tree-root> -> PASS(1/1), no SKIPPED block
```

## Steps

1. Create temp tree with only unlabeled leaves.
2. Run `doctest test <tree-root>`.

```go
func Setup(t *testing.T, req *Request) error {
	root := writeUnlabeledTree(t)
	req.Args = []string{"test", root}
	return nil
}
```