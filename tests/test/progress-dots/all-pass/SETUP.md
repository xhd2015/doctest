## Preconditions
- A temporary test tree with 3 passing leaves is created.

## Steps
1. Create test tree with 3 passing leaves.
2. Run `doctest test <dir>` (non-verbose).

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createPassFailTree(t, 3, 0)
    req.Args = []string{"test", testDir}
    return nil
}
```
