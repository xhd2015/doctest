## Preconditions
- The test command requires a doc-style test directory.

## Steps
1. Choose test arguments.
2. Run `doctest test`.

```go
import (
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 20 * time.Second
    return nil
}
```
