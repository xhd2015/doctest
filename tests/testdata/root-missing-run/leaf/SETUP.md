## Steps
1. Set the Input field on the request.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Input = "hello"
    return nil
}
```
