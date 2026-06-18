# Scenario

**Feature**: nested module path IS a child of parent module path

## Preconditions
- Nested `module` line is `ancestor/child` (prefix match).

## Steps
1. Create project with child module path.
2. Verify normal walk behavior is unchanged.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    t.Logf("child-path group")
    return nil
}
```
