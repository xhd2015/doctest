## Preconditions
- The root Setup has built the doctest binary and set a generous timeout.
- Child leaves test the `--gen-dir` flag behavior.

## Steps
1. Ensure the timeout from the parent Setup is preserved (no override needed at this grouping level).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
