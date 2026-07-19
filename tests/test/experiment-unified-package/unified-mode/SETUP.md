# Scenario

**Feature**: unified package generation + single-suite go test for a 2-leaf tree

```
# fixture: root DOCTEST defines ExperimentUnifiedRootMarker + Run; leaves a/, b/
RunTest(fixture, GenDir=tmp, ExperimentUnifiedPackagePerDoctestTree=true)
  -> go test ./suite passes both leaves
  -> gen has __droot, __registry, __allleaves, suite, leaf RunTestLeaf
```

## Preconditions

- Uses package suite entry `runner.RunTest` with explicit `GenDir` (`t.TempDir()`).
- Fixture built by root helper `createTwoLeafMarkerTree` unless a leaf overrides `Dir`.
- Only the unified Options field is set true; production must force ref.
- Layout metrics filled into `Response` after the run.

## Steps

1. Set `Op=run_gen` and `ExperimentUnifiedPackagePerDoctestTree=true`.
2. Run suite; assert pass and/or gen layout / package line / stderr.

## Context

- Sibling leaves under this node are MECE by assertion focus (pass vs layout vs suite-only vs announce).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "run_gen"
	req.ExperimentUnifiedPackagePerDoctestTree = true
	return nil
}
```
