# Scenario

**Feature**: P1 ref vs classic generation for a simple 2-leaf tree

```
# fixture: root DOCTEST defines ExperimentP1RootMarker + Run; leaves a/, b/
RunTest(fixture, GenDir=tmp, ExperimentRefInsteadOfInline=on|off)
  -> go test passes
  -> walk GenDir *.go for marker / type Request / imports
```

## Preconditions

- Uses package suite entry `runner.RunTest` with explicit `GenDir` (`t.TempDir()`).
- Fixture built by root helper `createTwoLeafMarkerTree` unless a leaf overrides `Dir`.
- Layout metrics filled into `Response` after the run (`MarkerDefCount`, leaf flags).

## Steps

1. Set `Op=ref_gen` and `ExperimentRefInsteadOfInline`.
2. Run suite; assert pass and/or gen layout.

## Context

- Flag off branch: classic inline (marker helper duplicated per leaf).
- Flag on branch: shared root package + thin leaf tests (P1).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "ref_gen"
	return nil
}
```
