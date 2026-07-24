# Scenario

**Feature**: key is deterministic for fixed inputs

```
# same spine + local pkgs + go version
ComputeLeafKey(in) -> k1
ComputeLeafKey(in) -> k2
# k1 == k2
```

## Preconditions

- Fixture is the base module/tree (no mutations).
- GoVersion is fixed at `go1.25.0`.

## Steps

1. Inherit base workspace from `key/`.
2. Set Op to `compute_twice`.

## Context

- Stability is the baseline for every sensitivity leaf below.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "compute_twice"
	return nil
}
```
