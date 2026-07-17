# Scenario

**Feature**: metrics-on suite records leaf_start / leaf_end for executed leaves

```
# 1-leaf pass fixture
RunTest(...) -> JSONL includes leaf_start and leaf_end (result pass)
```

## Preconditions

- Same tiny pass tree (one runnable leaf).
- Metrics enabled under injectable MetricsRoot.

## Steps

1. Run suite with metrics on.
2. Require leaf_start and leaf_end among decoded events.

## Context

- Discovery-only skips may be leaf_end skip without start; this leaf expects a
  real executed leaf so both start and end should appear when go-test JSON
  instrumentation is wired.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.NoMetrics = false
	return nil
}
```
