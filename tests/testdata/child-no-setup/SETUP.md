# Scenario

**Feature**: the root SETUP.md defines Request, Response, and Run

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Child No Setup Fixture

## Preconditions
- The root SETUP.md defines Request, Response, and Run.
- The leaf SETUP.md defines neither Setup nor Run, which triggers a validation
  error because non-root SETUP.md must have func Setup.

```go
import "fmt"

type Request struct {
    Value int
}

type Response struct {
    Result int
}

func Run(t *testing.T, req *Request) (*Response, error) {
    if req.Value < 0 {
        return nil, fmt.Errorf("Value must be non-negative")
    }
    return &Response{Result: req.Value * 2}, nil
}
```
