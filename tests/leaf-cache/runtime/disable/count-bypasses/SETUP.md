# Scenario

**Feature**: `-count=1` disables programmatic leaf-cache skip after a warm hit

```
run1: test fixture            -> store pass
run2: test fixture            -> Cached > 0
run3: test fixture -count=1   -> 0 Cached
```

## Preconditions

- Parent prepared fixture and warm Args/Args2.
- Any `-count=N` disables skip (using `1` as the canonical case).

## Steps

1. Set Args3 to `test <fixture> -count=1`.
2. Assert run2 Cached > 0; run3 Cached == 0; all exits 0.

## Context

- Aligns leaf-cache with go's "count disables cache" mental model.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Args3 = []string{"test", req.FixtureDir, "-count=1"}
	return nil
}
```
