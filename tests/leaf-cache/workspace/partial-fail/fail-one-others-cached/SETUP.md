# Scenario

**Feature**: after a partial workspace fail, the previously passing leaf is Cached on run2

```
run1: fail suite -> store only pass leaf
run2: Cached >= 1, still exit != 0 (fail leaf re-runs)
```

## Preconditions

- Parent prepared mixed pass/fail multi-tree fixture.

## Steps

1. Keep default double-run Args.
2. Assert both non-zero exits and run2 Cached >= 1.

## Context

- Fail leaf must never be skipped as warm pass.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	return nil
}
```
