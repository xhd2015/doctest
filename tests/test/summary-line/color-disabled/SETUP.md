# Scenario

**Feature**: --no-color produces plain PASS summary without ANSI

```
# color suppressed
doctest test --no-color <dir> -> plain PASS(passed/total)
```

## Preconditions

- A temp tree with 1 passing leaf exists.

## Steps

1. Create a 1-pass temp tree.
2. Run `doctest test --no-color <dir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	testDir := createPassFailTree(t, 1, 0)
	req.Args = []string{"test", "--no-color", testDir}
	return nil
}
```