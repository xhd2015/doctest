# Scenario

**Feature**: discovery test grouping

## Preconditions
- Parent `tests/SETUP.md` provides shared CLI runner types when running from `tests/` root.
- Child trees under `discovery/` are self-contained DOCTEST.md roots.

## Steps
1. Delegate to child feature trees.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    t.Logf("discovery group")
    return nil
}
```
