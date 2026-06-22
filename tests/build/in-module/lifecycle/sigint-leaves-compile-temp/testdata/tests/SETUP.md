# Setup

```go
import (
    "testing"
    "example.com/app/internal/greet"
)

type Request struct{}
type Response struct{ Message string }

func Run(t *testing.T, req *Request) (*Response, error) {
    return &Response{Message: greet.Hello()}, nil
}
```
