# Scenario

**Feature**: ancestor and nested module have different git contexts

## Preconditions
- `gitRoot(ancestor) != gitRoot(nested)` including mixed null/non-null cases.

## Steps
1. Create non-child module project with mismatched git contexts.
2. Verify warning emitted and nested module skipped.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    t.Logf("different-git group")
    return nil
}
```
