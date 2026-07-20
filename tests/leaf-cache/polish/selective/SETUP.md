# Scenario

**Feature**: invalidation is per-leaf / per-DAG, not whole-tree wipe

```
# sibling path
warm both leaves -> edit leaf_a ASSERT -> sibling still Cached

# local-dep path
warm leaf importing pkg -> edit pkg -> leaf re-runs (0 Cached)
```

## Preconditions

- Children choose sibling vs local-dep fixture (MECE).

## Steps

1. Child builds fixture and 3-run Args sequence with MutateAfterRun=2.

## Context

- Significance: selective miss is the main correctness property for multi-leaf suites.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	return nil
}
```
