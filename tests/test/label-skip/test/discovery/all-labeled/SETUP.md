# Scenario

**Feature**: all-labeled tree exits 0 with skip summary only

```
# only labeled leaf
doctest test <tree-root> -> 0 run, SKIPPED summary, no PASS/FAIL line
```

## Steps

1. Create labeled-only temp tree.
2. Run `doctest test <tree-root>`.

```go
func Setup(t *testing.T, req *Request) error {
	root := writeLabeledTree(t, false, "human-guided-ui-test", "manual only")
	req.Args = []string{"test", root}
	return nil
}
```