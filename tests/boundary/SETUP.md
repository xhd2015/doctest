# Scenario

**Feature**: this is the parent root of the boundary test tree

```
# DOCTEST.md creates an inheritance firewall
parent SETUP.md -/-> nested DOCTEST.md (no cross-inheritance)

# each root has its own Run, Request, Response, setup chain
nested root -> self-contained test tree -> runs independently
```

## Preconditions
- This is the parent root of the boundary test tree.

## Steps
1. The leaf calls Run which returns an error (stub, not implemented).

```go
import (
    "fmt"
    "testing"
)

type Request struct{}

type Response struct{}

func Run(t *testing.T, req *Request) (*Response, error) {
    return nil, fmt.Errorf("stub: not implemented")
}
```
