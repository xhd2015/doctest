# Scenario

**Feature**: default hierarchical ref packages under unified suite generation

```
# mini run under default gen
RunTest(tiny tree) -> pass

# hierarchical: root package + thin leaves (unified suite path)
RunTest(2-leaf marker fixture, GenDir=tmp)
  -> inspect gen *.go for ExperimentP1RootMarker / type Request / imports
```

## Preconditions

- Default generation is hierarchical unified (ref packages + suite).
- Root package owns `Request`, `Response`, `Run`, and root helpers such as
  `ExperimentP1RootMarker`.
- Leaf packages import root; do not redefine root types/helpers.
- Fixtures use distinctive `ExperimentP1RootMarker` /
  `ROOT_RUN_MARKER_P1_EXPERIMENT_REF` so layout can be counted without
  depending on exact root package dirname (`__droot`, …).
- Leaves never assert complex multi-level Setup edge cases.

## Steps

1. Leaf Setup sets `req.Op` and branch fields.
2. Root `Run` dispatches mini suite run or ref_gen (+ layout fill).
3. Leaf Assert checks run success or gen-layout metrics.

## Context

- Experiment CLI flags are removed; hierarchical unified is the only suite mode.
- Help token coverage lives in `tests/help/test-options` (parent CLI tree).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Leaves set Op and scenario fields; default keeps Run from erroring if forgotten.
	if req.Op == "" {
		req.Op = "mini_run"
	}
	return nil
}
```
