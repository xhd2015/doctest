## Steps
1. Override the Action field.
2. Define Run again — this redefines Run and is the violation under test.

```go
import (
    "fmt"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Action = "leaf-action"
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    return &Response{Result: fmt.Sprintf("child:%s", req.Action)}, nil
}
```
