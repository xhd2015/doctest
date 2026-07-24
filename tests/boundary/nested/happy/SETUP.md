# Scenario

**Feature**: the nested root defines Request{Name} and Response{Greeting}

```
# DOCTEST.md creates an inheritance firewall
parent SETUP.md -/-> nested DOCTEST.md (no cross-inheritance)

# each root has its own Run, Request, Response, setup chain
nested root -> self-contained test tree -> runs independently
```

## Preconditions
- The nested root defines Request{Name} and Response{Greeting}.
- Run returns a greeting when Name is provided.

## Steps
1. Set `req.Name` to `"World"`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Name = "World"
    return nil
}
```
