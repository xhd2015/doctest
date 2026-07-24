# Scenario

**Feature**: multiple directory arguments use one suite plan — one progress summary + one final aggregate

```
# multi-arg one suite plan (non-conflicting roots)
doctest test <dir-a> <dir-b> -> one prepare/workspace hub go test -> one (3 Run, 3 Pass, ...) in DURATION -> single PASS(3/3) in DURATION
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
