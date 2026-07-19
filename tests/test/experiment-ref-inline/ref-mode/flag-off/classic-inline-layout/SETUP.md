# Scenario

**Feature**: classic gen duplicates root marker helper into each leaf package

```
# flag off
gen/a/*_test.go defines ExperimentP1RootMarker
gen/b/*_test.go defines ExperimentP1RootMarker
# both leaves pass
```

## Preconditions

- Two-leaf marker fixture (default from `ref_gen`).

## Steps

1. Run with flag off and inspect GenDir.
2. Expect suite success and marker helper defined in **both** leaf test files
   (count ≥ 2 across gen tree).

## Context

- Does not require absence of every non-leaf path; requires classic duplication signal.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Reaffirm classic path; parent sets Op=ref_gen.
	req.ExperimentRefInsteadOfInline = false
	if req.Op != "ref_gen" {
		req.Op = "ref_gen"
	}
	return nil
}
```

