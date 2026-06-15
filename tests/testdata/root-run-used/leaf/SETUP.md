## Steps
1. Set the request Name to "doctest".

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Name = "doctest"
    return nil
}
```
