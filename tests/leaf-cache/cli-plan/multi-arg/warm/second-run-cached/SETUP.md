# Scenario

**Feature**: second default multi-arg run reports leaf-cache Cached for both trees

```
run1: multi-arg two trees -> exit 0, may be 0 Cached (cold)
run2: multi-arg two trees -> exit 0, sum(Cached) >= 2
```

## Preconditions

- Parent prepared all-pass multi-tree fixture and identical multi-arg Args/Args2.

## Steps

1. Keep double-run configuration.
2. Assert run2 total Cached >= 2 and both exits 0.

## Context

- Proves PutPass + warm GetPass skip on multi-arg path for both roots.
- Engine-agnostic: works if multi-arg fans to N× TestWithStats (two summary
  lines with 1 Cached each) or a unified plan (one line with 2 Cached).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	return nil
}
```
