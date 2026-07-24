# Scenario

**Feature**: root marker helper defined once under gen (shared __droot)

```
RunTest(2-leaf, GenDir=tmp)
  -> exactly one func ExperimentP1RootMarker definition
```

## Preconditions

- Default hierarchical generation.
- GenDir kept for walk.

## Steps

1. `Op=ref_gen`.
2. Assert MarkerDefCount == 1.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "ref_gen"
	return nil
}
```
