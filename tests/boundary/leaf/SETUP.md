# Scenario

**Feature**: the parent root defines Request{} and Response{} with a stub Run that returns an error

```
# DOCTEST.md creates an inheritance firewall
parent SETUP.md -/-> nested DOCTEST.md (no cross-inheritance)

# each root has its own Run, Request, Response, setup chain
nested root -> self-contained test tree -> runs independently
```

## Preconditions
- The parent root defines Request{} and Response{} with a stub Run that returns an error.

## Steps
1. No additional setup needed; Run is called automatically with an empty Request.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    t.Log("boundary: parent root leaf, stub Run will return error")
    return nil
}
```
