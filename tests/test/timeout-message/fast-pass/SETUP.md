# Scenario

**Feature**: fast-passing suites do not emit a false timeout Error

```
# fast 1-pass tree under normal/generous timeout
doctest test --no-color <pass-tree> -> go test completes -> exit 0

# must not look like a timeout failure
doctest -> no "Error: go test timed out"
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Fixture is a single fast-passing leaf (no sleep).

## Steps

1. Create a temp 1-pass tree.
2. Run `doctest test --no-color <tree>` (default timeout policy; no short --timeout).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	testDir := createFastPassTree(t)
	req.Args = []string{"test", "--no-color", testDir}
	return nil
}
```
