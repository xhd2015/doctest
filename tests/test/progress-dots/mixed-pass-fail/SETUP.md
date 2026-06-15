## Preconditions
- A temporary test tree with 2 passing and 1 failing leaves is created.

## Steps
1. Create test tree with 2 pass + 1 fail.
2. Run `doctest test <dir>` (non-verbose).

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createPassFailTree(t, 2, 1)
    req.Args = []string{"test", testDir}
    return nil
}
```
