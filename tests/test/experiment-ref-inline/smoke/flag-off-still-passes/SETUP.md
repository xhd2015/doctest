# Scenario

**Feature**: default path (flag off) still runs a tiny tree

```
RunTest(1-leaf, ExperimentRefInsteadOfInline=false) -> success
```

## Preconditions

- Field left at default false (explicit in Request for clarity).

## Steps

1. Run one-leaf pass fixture with flag off.
2. Expect no run error.

## Context

- Regression guard: introducing the Options field must not break classic runs.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.ExperimentRefInsteadOfInline = false
	return nil
}
```
