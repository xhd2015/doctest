# Scenario

**Feature**: root helper appears once under gen when flag is on

```
# flag on
count(files defining func ExperimentP1RootMarker) == 1
# (classic would be ≥2)
```

## Preconditions

- GenDir retained after run for layout walk.

## Steps

1. Run with flag on.
2. Assert suite success and `MarkerDefCount == 1`.

## Context

- Root package dirname is not fixed (`__droot` / `_root` / …); count by symbol.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.ExperimentRefInsteadOfInline = true
	req.Op = "ref_gen"
	return nil
}
```

