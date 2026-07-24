# Scenario

**Feature**: ancestor and nested module share the same git work tree

## Preconditions
- `gitRoot(ancestor) == gitRoot(nested)` and both are non-null (or both null in sibling branch).

## Steps
1. Create non-child module project in a single git repo.
2. Verify nested module doctest trees are discovered.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    t.Logf("same-git group")
    return nil
}
```
