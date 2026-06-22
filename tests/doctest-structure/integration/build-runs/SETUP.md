# Scenario

**Feature**: `doctest build` succeeds on a minimal new-layout tree

```
# build assembles from DOCTEST.md types
temp tree (new layout + leaf) -> doctest build -> exit 0
```

## Preconditions

- A minimal valid tree with version, DSN, types in `DOCTEST.md`, and one leaf.

## Steps

1. Write the minimal valid tree with a runnable leaf.
2. Run `doctest build <treeDir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	treeDir := writeTree(t, treeOpts{withVersion: true, withLeaf: true})
	req.Args = []string{"build", treeDir}
	return nil
}
```