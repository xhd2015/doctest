# Scenario

**Feature**: full multi-arg generate then subset keeps sibling manifest entries

```
run1: test tree-a tree-b --gen-dir G
run2: test tree-a --gen-dir G
  -> manifest still has tree-b/*
  -> G/tree-b still on disk
```

## Preconditions

- Root Setup provided product binary.
- Isolated two-tree fixture + GenDir.

## Steps

1. prepareTwoTreeModule.
2. ArgsFull = both trees; ArgsSubset = tree-a only.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareTwoTreeModule(t, req)
	// Relative multi-arg from WorkDir (same as product CLI usage).
	req.ArgsFull = baseArgs(req, "tree-a", "tree-b")
	req.ArgsSubset = baseArgs(req, "tree-a")
	return nil
}
```
