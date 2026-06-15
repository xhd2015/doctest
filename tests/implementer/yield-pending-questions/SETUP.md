## Preconditions
- Tests in this group run the `yield-pending-questions` binary directly.
- The binary is dispatched via the doctest binary copied as `yield-pending-questions`.

## Steps
1. Read `YIELD_PQ_BIN` from environment.
2. Set the binary for execution.

```go
import (
    "os"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=yield-pending-questions")
    req.Bin = os.Getenv("YIELD_PQ_BIN")
    return nil
}
```
