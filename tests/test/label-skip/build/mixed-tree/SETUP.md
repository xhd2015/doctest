# Scenario

**Feature**: build compiles mixed fast and labeled leaves

```
# fast + labeled in one tree
doctest build <tree> -> exit 0, compiles all leaves
```

## Steps

1. Create temp tree with fast_leaf and labeled_leaf.
2. Run `doctest build <tree>`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	root := writeLabeledTree(t, true, "ui-automation", "mixed build")
	req.Args = []string{"build", root}
	return nil
}
```