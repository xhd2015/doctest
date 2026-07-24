# Scenario

**Feature**: three passing leaves produce an all-pass summary

```
# discover leaves, run packages, aggregate stats
doctest test <dir> -> discover leaves -> go test per leaf -> accumulate Passed/Total

# final summary
runner -> PASS(passed/total)
```

## Preconditions

- A temp tree with 3 passing leaves exists.

## Steps

1. Create a 3-pass / 0-fail temp tree.
2. Run `doctest test <dir>` (non-verbose).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createPassFailTree(t, 3, 0)
	req.Args = []string{"test", testDir}
	return nil
}
```