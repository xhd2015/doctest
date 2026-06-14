## Preconditions
- This is a nested root with its own DOCTEST.md boundary.
- This root defines its own Request and Response types (different from parent root).

## Steps
1. The leaf sets `req.Name`.
2. Run returns a greeting using the Name field, or an error if Name is empty.

```go
import (
    "fmt"
    "testing"
)

type Request struct {
    Name string
}

type Response struct {
    Greeting string
}

func Run(t *testing.T, req *Request) (*Response, error) {
    if req.Name == "" {
        return nil, fmt.Errorf("Name is required")
    }
    return &Response{Greeting: "Hello, " + req.Name + "!"}, nil
}
```
