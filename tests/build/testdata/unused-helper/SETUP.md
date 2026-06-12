## Setup

- Defines three helper functions (helperA, helperB, helperC).
- Only helperA is used by func Run.
- The generator emits all helpers as closures; unused helperB and helperC
  cause "declared and not used" compilation errors.

```go
import "fmt"

type Request struct {
    Name string
}

type Response struct {
    Message string
}

func helperA(s string) string {
    return "A: " + s
}

func helperB(s string) string {
    return "B: " + s
}

func helperC(s string) string {
    return "C: " + s
}

func Setup(t *testing.T, req *Request) error {
    req.Name = "default"
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    msg := helperA(req.Name)
    _ = fmt.Sprintf("%v", msg)
    return &Response{Message: msg}, nil
}
```
