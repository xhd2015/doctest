# Scenario

**Feature**: build compiles labeled-only tree

```
# labeled leaf only
doctest build <tree> -> compiles generated Go, exit 0
```

## Steps

1. Create labeled-only temp tree.
2. Run `doctest build <tree>`.

```go
func Setup(t *testing.T, req *Request) error {
	root := writeLabeledTree(t, false, "ui-automation", "compile check")
	req.Args = []string{"build", root}
	return nil
}
```