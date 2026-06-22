# Scenario

**Feature**: mixed pass and fail leaves produce FAIL summary after failure output

```
# failure detail before summary
go test -> FAIL\t lines -> runner -> FAIL(passed/total)
```

## Preconditions

- A temp tree with 2 passing and 1 failing leaf exists.

## Steps

1. Create a 2-pass / 1-fail temp tree.
2. Run `doctest test <dir>` (non-verbose).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	testDir := createPassFailTree(t, 2, 1)
	req.Args = []string{"test", testDir}
	return nil
}
```