# Scenario

**Feature**: default hierarchical unified package-per-tree generation

```
# default path (no experiment flags)
RunTest(2-leaf marker fixture, GenDir=tmp)
  -> __droot + __registry + leaf RunTestLeaf + __allleaves + suite
  -> go test suite only; both leaves pass
```

## Preconditions

- Default generation is hierarchical unified (ref packages + suite).
- Fixtures use distinctive `ExperimentUnifiedRootMarker` /
  `ROOT_RUN_MARKER_UNIFIED_PACKAGE` so layout can be counted without
  hard-coding package dirnames.
- Leaves never assert complex multi-level Setup edge cases.

## Steps

1. Leaf Setup sets `req.Op=run_gen` and optional Dir/GenDir.
2. Root `Run` builds fixture (if needed), runs `runner.RunTest` with default Options, fills layout.
3. Leaf Assert checks run success, gen-layout, or suite-only go test packaging.

## Context

- Help no longer documents experiment flags (`tests/help/test-options`).
- Sibling tree `tests/test/experiment-ref-inline/` asserts hierarchical ref properties under the same default.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if req.Op == "" {
		req.Op = "run_gen"
	}
	return nil
}
```
