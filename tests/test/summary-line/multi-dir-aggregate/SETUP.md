# Scenario

**Feature**: multiple directory arguments produce one aggregated summary

```
# multi-arg aggregation
doctest test <dir-a> <dir-b> -> accumulate stats -> single PASS(passed/total)
```

## Preconditions

- Two temp trees exist: one with 2 passing leaves, one with 1 passing leaf.

## Steps

1. Create tree A (2 pass) and tree B (1 pass).
2. Run `doctest test <tree-a> <tree-b>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dirA := createPassFailTree(t, 2, 0)
	dirB := createPassFailTree(t, 1, 0)
	req.Args = []string{"test", dirA, dirB}
	return nil
}
```