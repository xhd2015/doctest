# Scenario

**Feature**: spine Go content is part of the leaf key

```
# before
spine(v1) -> key1

# edit leaf ASSERT Go or ancestor SETUP Go
spine(v2) -> key2
# key1 != key2
```

## Preconditions

- Base fixture; mutation targets only doctest markdown Go blocks on the spine.
- Op is `compute_mutate`.

## Steps

1. Compute key on the original spine.
2. Apply leaf or ancestor mutation.
3. Recompute; keys must differ.

## Context

- Leaf assert and ancestor setup are MECE spine layers; root DOCTEST is covered
  implicitly by fixture presence (stable leaf) but not re-mutated here.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "compute_mutate"
	return nil
}
```
