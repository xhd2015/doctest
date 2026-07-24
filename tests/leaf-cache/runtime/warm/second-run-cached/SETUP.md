# Scenario

**Feature**: second default run of a previously passing leaf reports Cached > 0

```
run1: doctest test fixture -> exit 0, may be 0 Cached (cold)
run2: doctest test fixture -> exit 0, Cached >= 1
```

## Preconditions

- Parent prepared 1-pass fixture and identical Args/Args2.
- Isolated DOCTEST_LEAF_CACHE; fresh GOCACHE per run from runtime root.

## Steps

1. Keep default double-run configuration.
2. Assert run2 summary Cached count > 0.

## Context

- Proves PutPass on success + GetPass skip + summary integration.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	// Parent already set fixture and args; reaffirm Op.
	req.Op = "runtime_multi"
	return nil
}
```
