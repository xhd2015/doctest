## Preconditions
- The nested root defines Request{Name} and Response{Greeting}.
- Run returns a greeting when Name is provided.

## Steps
1. Set `req.Name` to `"World"`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Name = "World"
    return nil
}
```
