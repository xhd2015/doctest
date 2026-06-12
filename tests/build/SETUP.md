## Preconditions
- The build command requires a doc-style test directory.

## Steps
1. Choose build arguments.
2. Run `doctest build`.

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
