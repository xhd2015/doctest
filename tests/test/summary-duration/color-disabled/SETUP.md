# Scenario

**Feature**: --no-color produces plain durations without ANSI

```
# color suppressed
doctest test --no-color <dir> -> plain inline DURATION -> plain PASS(p/t) in DURATION
```

## Preconditions

- A temp tree with 1 passing leaf exists.

## Steps

1. Create a 1-pass temp tree.
2. Run `doctest test --no-color <dir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	testDir := createPassFailTree(t, 1, 0)
	req.Args = []string{"test", "--no-color", testDir}
	return nil
}
```