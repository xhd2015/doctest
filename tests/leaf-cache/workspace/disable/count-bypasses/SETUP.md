# Scenario

**Feature**: workspace `-count=1` yields 0 Cached after a proven warm hit

```
run1: store
run2: Cached > 0
run3: -count=1 -> 0 Cached
```

## Preconditions

- Parent prepared fixture and warm Args/Args2.

## Steps

1. Set Args3 to `test <mod>/... -count=1`.
2. Assert run2 Cached > 0; run3 Cached == 0; all exits 0.

## Context

- Any `-count=N` disables skip (using `1` as the canonical case).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Args3 = []string{"test", mustWorkspacePattern(req.WorkDir), "-count=1"}
	return nil
}
```
