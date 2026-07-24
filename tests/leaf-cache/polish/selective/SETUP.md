# Scenario

**Feature**: selective leaf-key invalidation (L2 library)

```
# sibling: mutate leaf_a ASSERT → leaf_b key stable, leaf_a key changes
# local dep: mutate imported pkg → leaf key changes
```

## Preconditions

- In-process `ComputeLeafKey` only (no product binary).
- Unlabeled.

## Steps

1. Child builds fixture and sets Op/Mutation.
2. Assert key stability / change flags.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	return nil
}
```
