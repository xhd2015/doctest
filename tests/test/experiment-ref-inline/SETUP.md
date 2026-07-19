# Scenario

**Feature**: experiment ref-instead-of-inline — P0 Options plumbing + P1 package DAG

```
# P0: parse flag into Options (default off)
parseTestOptions([...]) -> Options.ExperimentRefInsteadOfInline

# P0/P1 mini run
RunTest(tiny tree, ExperimentRefInsteadOfInline=?) -> pass

# P1: flag off → classic per-leaf inline; flag on → root package + thin leaves
RunTest(2-leaf marker fixture, GenDir=tmp, ExperimentRefInsteadOfInline=?)
  -> inspect gen *.go for ExperimentP1RootMarker / type Request / imports
```

## Preconditions

- P0 symbols (implemented):
  - `core.Options.ExperimentRefInsteadOfInline bool`
  - `runner.ParseTestOptions` understands `--experiment-ref-instead-of-inline`
- P1 symbols/behavior expected (implementer; **RED** until present):
  - When `ExperimentRefInsteadOfInline` is true, generation uses a **ref**
    package DAG (not classic `AssembleTestSource` only)
  - Root package owns `Request`, `Response`, `Run`, and root helpers such as
    `ExperimentP1RootMarker`
  - Leaf `_test.go` files import root; do not redefine root types/helpers
  - Flag false remains classic: each leaf package inlines ancestor source
  - Prefer gen-cache isolation from warm classic mapping-gen (explicit
    `GenDir` in tests; production may use `mapping-gen-ref` or a mode marker)
- P1 fixtures use distinctive `ExperimentP1RootMarker` /
  `ROOT_RUN_MARKER_P1_EXPERIMENT_REF` so layout can be counted without
  depending on exact root package dirname (`__droot`, `_root`, …).
- Leaves never assert complex multi-level Setup edge cases (P2).

## Steps

1. Leaf Setup sets `req.Op` and branch fields.
2. Root `Run` dispatches parse, mini suite run, or ref_gen (+ layout fill).
3. Leaf Assert checks option bits, run success, or gen-layout metrics.

## Context

- Hard product rule: **without the flag, behavior must not change** (classic only).
- Help token coverage lives in `tests/help/test-options` (parent CLI tree).
- Sealed P0 leaves under `flags/` and `smoke/` must keep passing once P0 is green;
  new P1 leaves live only under `ref-mode/`.

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
