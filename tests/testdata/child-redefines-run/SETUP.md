# Scenario

**Feature**: the root SETUP.md defines Request, Response, and Run

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Child Redefines Run Fixture

## Preconditions
- The root SETUP.md defines Request, Response, and Run.
- A child SETUP.md also defines Run, which should trigger a validation error.

```go
import "fmt"

type Request struct {
    Action string
}

type Response struct {
    Result string
}

func Run(t *testing.T, req *Request) (*Response, error) {
    if req.Action == "" {
        return nil, fmt.Errorf("Action is required")
    }
    return &Response{Result: "root:" + req.Action}, nil
}
```
