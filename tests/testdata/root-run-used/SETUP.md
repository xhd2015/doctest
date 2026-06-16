# Scenario

**Feature**: this fixture verifies that only the root Run is used in generated code,

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Root Run Used Fixture

## Preconditions
- This fixture verifies that only the root Run is used in generated code,
  even when a child is in the ancestor chain.

## Steps
1. The root Run returns a greeting message including the request Name.
2. The leaf Setup sets the Name.
3. The leaf Assert checks the message comes from the root Run.

```go
import "fmt"

type Request struct {
    Name string
}

type Response struct {
    Message string
}

func Setup(t *testing.T, req *Request) error {
    if req.Name == "" {
        req.Name = "world"
    }
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    return &Response{Message: fmt.Sprintf("hello %s", req.Name)}, nil
}
```
