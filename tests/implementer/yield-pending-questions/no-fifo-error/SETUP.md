## Preconditions
- `QUESTION_FIFO` environment variable is NOT set.

## Steps
1. Invoke yield-pending-questions without `QUESTION_FIFO`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{`{"id":"1","question":"test"}`}
    return nil
}
```
