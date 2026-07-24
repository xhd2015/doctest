# Scenario

**Feature**: `doctest test` succeeds on a minimal new-layout tree

```
# test discovers and runs leaves
temp tree (new layout + leaf) -> doctest test -> leaf passes, exit 0
```

## Preconditions

- A minimal valid tree with version, DSN, types in `DOCTEST.md`, and one passing leaf.

## Steps

1. Write the minimal valid tree with a runnable leaf.
2. Run `doctest test <treeDir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	treeDir := writeTree(t, treeOpts{withVersion: true, withLeaf: true})
	req.Args = []string{"test", treeDir}
	return nil
}
```