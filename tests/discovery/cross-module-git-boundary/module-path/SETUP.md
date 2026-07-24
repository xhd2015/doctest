# Scenario

**Feature**: module path relationship at go.mod boundary

## Preconditions
- Parent and nested modules differ in whether the nested `module` path is a child of the ancestor path.

## Steps
1. Branch on child-path vs non-child-path module relationship.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    t.Logf("module-path group")
    return nil
}
```
