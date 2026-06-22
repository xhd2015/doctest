# Scenario

**Feature**: multiple directory arguments produce per-dir inline durations and one final aggregate

```
# multi-arg aggregation
doctest test <dir-a> <dir-b> -> per-dir inline (N Run, ...) in DURATION -> single PASS(passed/total) in DURATION
```

## Preconditions

- Two temp trees exist: one with 2 passing leaves, one with 1 passing leaf.

## Steps

1. Create tree A (2 pass) and tree B (1 pass).
2. Run `doctest test <tree-a> <tree-b>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	dirA := createPassFailTree(t, 2, 0)
	dirB := createPassFailTree(t, 1, 0)
	req.Args = []string{"test", dirA, dirB}
	return nil
}
```