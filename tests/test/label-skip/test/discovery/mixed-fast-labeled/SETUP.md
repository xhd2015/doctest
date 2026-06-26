# Scenario

**Feature**: mixed tree runs fast leaf and skips labeled leaf

```
# one fast + one labeled under tree root
doctest test <tree-root> -> PASS(1/1) + SKIPPED summary
```

## Steps

1. Create temp tree with `fast_leaf` and `labeled_leaf`.
2. Run `doctest test <tree-root>`.

```go
func Setup(t *testing.T, req *Request) error {
	root := writeLabeledTree(t, true, "ui-automation", "heavy ui test")
	req.Args = []string{"test", root}
	return nil
}
```