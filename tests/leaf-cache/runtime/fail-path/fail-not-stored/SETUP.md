# Scenario

**Feature**: two runs of a failing leaf never report Cached and never skip

```
run1: fail -> exit != 0, 0 Cached
run2: fail -> exit != 0, 0 Cached
```

## Preconditions

- Parent fail fixture and default Args/Args2.

## Steps

1. Keep double-run fail configuration.
2. Assert both exits non-zero and Cached == 0.

## Context

- If run2 showed Cached > 0, a fail was incorrectly PutPass'd.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	return nil
}
```
