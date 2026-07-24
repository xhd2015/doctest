# Scenario

**Feature**: leaf_end elapsed is package-attributed, not full tree wall cloned to every leaf

```
# single-leaf: leaf elapsed ≈ go_test package time, not infinite
# multi-leaf would differ per package; 1-leaf still must have finite elapsed_ns
```

## Preconditions

- MetricsOn true.

## Steps

1. Run 1-leaf fixture with metrics.
2. Compare leaf_end.elapsed_ns to go_test phase elapsed_ns.

## Context

- Previously every leaf got the whole tree wall; single-leaf may still equal go_test.
- Assert leaf elapsed is positive and not larger than suite-ish bound (go_test phase * 2).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "record_run"
	req.MetricsOn = true
	return nil
}
```
