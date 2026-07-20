# Scenario

**Feature**: two ComputeLeafKey calls with identical inputs return the same hex key

```
# first call
in -> ComputeLeafKey -> key1

# second call, no disk or source changes
in -> ComputeLeafKey -> key2
# key1 == key2, both lowercase hex
```

## Preconditions

- Base fixture from ancestors (app module, group/leaf spine, helper package).
- No mutation between calls.

## Steps

1. Keep Op=`compute_twice` from parent.
2. Run twice via root Run.
3. Assert equal non-empty lowercase hex digests.

## Context

- If keys differ without mutation, the DAG hash is non-deterministic (forbidden).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	// Op already compute_twice; fixture already built.
	req.Op = "compute_twice"
	return nil
}
```
