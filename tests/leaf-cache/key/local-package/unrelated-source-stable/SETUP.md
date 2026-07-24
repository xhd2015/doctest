# Scenario

**Feature**: editing a local package not imported by the spine leaves the key stable

```
# unrelated/noise.go is not in the import closure
noise v1 -> key1
noise v2 -> key2
# key1 == key2
```

## Preconditions

- Base fixture includes `unrelated/noise.go` with no imports from spine.
- Mutation = `local_unrelated`.

## Steps

1. Set Mutation to `local_unrelated`.
2. Run compute_mutate.
3. Assert keys equal.

## Context

- Hashing the entire module would over-invalidate; only the local import closure counts.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "compute_mutate"
	req.Mutation = "local_unrelated"
	return nil
}
```
