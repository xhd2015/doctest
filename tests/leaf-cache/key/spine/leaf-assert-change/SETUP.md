# Scenario

**Feature**: changing leaf ASSERT Go changes the key

```
# original leaf ASSERT Go
ASSERT.md Go block -> key1

# mutate assert body (comment marker in Go)
ASSERT.md Go block' -> key2
# key1 != key2
```

## Preconditions

- Base fixture leaf at `group/leaf`.
- Mutation = `leaf_assert`.

## Steps

1. Set Mutation to `leaf_assert`.
2. Run compute_mutate.
3. Assert Key ≠ Key2.

## Context

- Leaf assert is the most leaf-private spine layer; suite skip (P2) will use this.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "compute_mutate"
	req.Mutation = "leaf_assert"
	return nil
}
```
