# Scenario

**Feature**: multi-arg `-count=1` yields 0 Cached after a proven warm hit

```
run1: store
run2: Cached > 0
run3: -count=1 -> 0 Cached
```

## Preconditions

- Parent prepared fixture and warm multi-arg Args/Args2.

## Steps

1. Set Args3 to `test tree-a tree-b -count=1`.
2. Assert run2 total Cached >= 2; run3 total Cached == 0; all exits 0.

## Context

- Any `-count=N` disables skip (using `1` as the canonical case).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Args3 = []string{"test", req.TreeRoot, req.TreeRootB, "-count=1"}
	return nil
}
```
