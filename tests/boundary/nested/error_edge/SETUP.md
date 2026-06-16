# Scenario

**Feature**: this is a deeply nested root with its own DOCTEST.md boundary

```
# DOCTEST.md creates an inheritance firewall
parent SETUP.md -/-> nested DOCTEST.md (no cross-inheritance)

# each root has its own Run, Request, Response, setup chain
nested root -> self-contained test tree -> runs independently
```

## Preconditions
- This is a deeply nested root with its own DOCTEST.md boundary.
- Types are completely independent from ancestor roots.

## Steps
1. The leaf sets `req.ID` and `req.Data`.
2. Run validates and processes the request.

```go
import (
    "fmt"
    "testing"
)

type Request struct {
    ID   int
    Data string
}

type Response struct {
    Status  string
    Message string
}

func Run(t *testing.T, req *Request) (*Response, error) {
    if req.ID <= 0 {
        return nil, fmt.Errorf("ID must be positive")
    }
    if req.Data == "" {
        return nil, fmt.Errorf("Data is required")
    }
    return &Response{Status: "ok", Message: "processed: " + req.Data}, nil
}
```
