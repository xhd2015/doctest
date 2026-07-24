# Scenario

**Feature**: nested module path is NOT a child of parent module path

## Preconditions
- Nested `module` line does not have `ancestor/` prefix (sibling module).

## Steps
1. Branch on git context comparison at the module boundary.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    t.Logf("non-child-path group")
    return nil
}
```
