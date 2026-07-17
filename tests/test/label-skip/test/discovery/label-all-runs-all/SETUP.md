# Scenario

**Feature**: --label-all runs fast and labeled leaves under discovery

```
# mixed tree + --label-all
doctest test --label-all <tree-root> -> PASS(2/2), no SKIPPED for labels
```

## Steps

1. Create temp tree with `fast_leaf` and `labeled_leaf`.
2. Run `doctest test --label-all <tree-root>`.

```go
func Setup(t *testing.T, req *Request) error {
	root := writeLabeledTree(t, true, "ui-automation", "heavy ui test")
	req.Args = []string{"test", "--label-all", root}
	return nil
}
```
