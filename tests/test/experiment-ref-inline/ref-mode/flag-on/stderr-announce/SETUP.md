# Scenario

**Feature**: optional stderr announce when ref mode is active

```
RunTest(..., ExperimentRefInsteadOfInline=true)
  -> stderr mentions experiment ref-instead-of-inline (meta line)
```

## Preconditions

- Soft P1 criterion; wording may be gray meta similar to cold-cache announce.

## Steps

1. Run with flag on.
2. Assert suite success and stderr contains both `experiment` and
   `ref-instead-of-inline` (case-insensitive check in Assert).

## Context

- Does not require ANSI color tokens; substring only.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.ExperimentRefInsteadOfInline = true
	req.Op = "ref_gen"
	return nil
}
```

