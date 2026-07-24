# Scenario

**Feature**: a single passing leaf produces PASS(1/1)

```
# discover leaves, run packages, aggregate stats
doctest test <dir> -> discover 1 leaf -> PASS(1/1)
```

## Preconditions

- A temp tree with 1 passing leaf exists.

## Steps

1. Create a 1-pass / 0-fail temp tree.
2. Run `doctest test <dir>` (non-verbose).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createPassFailTree(t, 1, 0)
	req.Args = []string{"test", testDir}
	return nil
}
```