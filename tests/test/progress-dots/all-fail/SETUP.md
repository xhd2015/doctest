## Preconditions
- A temporary test tree with 3 failing leaves is created.

## Steps
1. Create test tree with 3 failing leaves.
2. Run `doctest test <dir>` (non-verbose).

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createPassFailTree(t, 0, 3)
    req.Args = []string{"test", testDir}
    return nil
}
```
