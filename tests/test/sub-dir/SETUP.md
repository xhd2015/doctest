## Preconditions
- The test command supports sub-directory resolution.

## Steps
1. Create a temporary doc-style test tree.
2. Run `doctest test <sub-dir>` on the tree.

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
