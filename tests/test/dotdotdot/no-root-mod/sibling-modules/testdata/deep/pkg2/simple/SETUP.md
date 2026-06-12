```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    t.Logf("setup: %s", req.Name)
    return nil
}
```
