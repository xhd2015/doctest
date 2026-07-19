# Scenario

**Feature**: enabled metrics write phase events for the tree pipeline

## Preconditions

- MetricsOn true (parent Setup).

## Steps

1. record_run on default 1-leaf tree.
2. Collect type=phase events.

## Context

- Expect at least discover, generate, go_test (materialize may be very short but still emitted).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Ensure record_run with metrics; parent already sets MetricsOn.
	req.Op = "record_run"
	req.MetricsOn = true
	return nil
}
```
