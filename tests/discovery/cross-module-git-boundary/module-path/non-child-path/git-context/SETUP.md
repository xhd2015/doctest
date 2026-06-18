# Scenario

**Feature**: git context at non-child module boundary

## Preconditions
- Non-child module path triggers `gitRoot(ancestor)` vs `gitRoot(nested)` comparison.

## Steps
1. Branch on equal (discover) vs unequal (warn + skip) git roots.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    t.Logf("git-context group")
    return nil
}
```
