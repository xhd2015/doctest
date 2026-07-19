# Scenario

**Feature**: leaf `_test.go` is thin — imports root, no inlined root types/helpers

```
# flag on
leaf a/b *_test.go:
  - does NOT define ExperimentP1RootMarker
  - does NOT contain "type Request"
  - DOES import a non-stdlib package (root package path)
```

## Preconditions

- Same 2-leaf marker fixture and GenDir walk as siblings.

## Steps

1. Run with flag on.
2. For each discovered leaf test file, assert thin import shape.

## Context

- Assert body stays leaf-local in the fixture; only root types/helpers move to the root package.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.ExperimentRefInsteadOfInline = true
	req.Op = "ref_gen"
	return nil
}
```

