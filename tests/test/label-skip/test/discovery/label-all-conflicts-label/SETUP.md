# Scenario

**Feature**: --label-all rejects combination with --label

```
doctest test --label-all --label heavy <tree> -> non-zero, mutual exclusion error
```

## Steps

1. Create any valid temp tree.
2. Run with both flags.

```go
func Setup(t *testing.T, req *Request) error {
	root := writeLabeledTree(t, true, "ui-automation", "heavy ui test")
	req.Args = []string{"test", "--label-all", "--label", "heavy", root}
	return nil
}
```
