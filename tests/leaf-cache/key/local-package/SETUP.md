# Scenario

**Feature**: local package import closure participates in the key

```
# spine imports example.com/app/pkg/helper
# unrelated/ is in the module but not imported

helper.go change -> key changes
noise.go change  -> key stable
```

## Preconditions

- Base fixture: `pkg/helper` imported by leaf ASSERT; `unrelated` is not.
- Op is `compute_mutate`.

## Steps

1. Child picks Mutation `local_imported` or `local_unrelated`.
2. Compare keys before/after mutation.

## Context

- Collectively exhaustive for local packages: on-closure vs off-closure.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "compute_mutate"
	return nil
}
```
