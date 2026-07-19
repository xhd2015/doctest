# Scenario

**Feature**: experiment unified-package-per-doctest-tree — one go test binary per tree

```
# Parse flag into Options (default off; on implies ref)
parseTestOptions([...]) -> Options.ExperimentUnifiedPackagePerDoctestTree
                        -> Options.ExperimentRefInsteadOfInline (forced when unified)

# Unified gen + run (flag on)
RunTest(2-leaf marker fixture, GenDir=tmp, ExperimentUnifiedPackagePerDoctestTree=true)
  -> gen: __droot + __registry + leaf RunTestLeaf + __allleaves + suite
  -> go test ./suite (one package)
  -> both leaves pass

# Control (flag off)
RunTest(same fixture, unified=false) -> classic multi-leaf *_test.go
```

## Preconditions

- Symbols/behavior expected (implementer; **RED** until present):
  - `core.Options.ExperimentUnifiedPackagePerDoctestTree bool`
  - `runner.ParseTestOptions` understands `--experiment-unified-package-per-doctest-tree`
  - When unified is true, `ExperimentRefInsteadOfInline` is forced true
  - Generation layout under tree gen root:
    - `__droot/`, `__registry/`, `__allleaves/`, leaf non-test with `RunTestLeaf`, `suite/suite_test.go`
  - Suite imports registry + allleaves (+ stdlib os for DOCTEST_METRICS_PARENT_LEAF)
  - `go test` only the suite package → one test binary per DOCTEST tree
  - Flag false remains classic: multi-package leaf `*_test.go`
- Fixtures use distinctive `ExperimentUnifiedRootMarker` /
  `ROOT_RUN_MARKER_UNIFIED_PACKAGE`.
- Help token coverage lives in `tests/help/test-options` (parent CLI tree).

## Steps

1. Leaf Setup sets `req.Op` and branch fields.
2. Root `Run` dispatches parse or run_gen (+ layout fill).
3. Leaf Assert checks option bits, run success, gen layout, or go-test package line.

## Context

- Hard product rule: **without the flag, behavior must not change** (classic only).
- Unified **auto-enables** ref; tests set only the unified Options field on RunTest.
- Sealed sibling tree `tests/test/experiment-ref-inline/` must keep passing.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Leaves set Op and scenario fields; default keeps Run from erroring if forgotten.
	if req.Op == "" {
		req.Op = "parse_flags"
	}
	return nil
}
```
