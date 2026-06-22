# Scenario

**Feature**: a single passing leaf reports duration in both summaries

```
# discover leaves, run packages, measure wall time
doctest test <dir> -> discover 1 leaf -> inline (1 Run, ...) in DURATION -> PASS(1/1) in DURATION
```

## Preconditions

- A temp tree with 1 passing leaf exists.

## Steps

1. Create a 1-pass / 0-fail temp tree.
2. Run `doctest test <dir>` (non-verbose).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	testDir := createPassFailTree(t, 1, 0)
	req.Args = []string{"test", testDir}
	return nil
}
```