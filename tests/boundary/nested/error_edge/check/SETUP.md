## Preconditions
- The deeply nested root defines Request{ID, Data} and Response{Status, Message}.
- Run validates ID > 0 and Data non-empty.

## Steps
1. Set `req.ID = 42` and `req.Data = "test"`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.ID = 42
    req.Data = "test"
    return nil
}
```
