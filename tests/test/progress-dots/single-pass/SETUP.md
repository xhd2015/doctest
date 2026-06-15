## Preconditions
- A temporary test tree with 1 passing leaf is created.

## Steps
1. Create test tree with 1 passing leaf.
2. Run `doctest test <dir>` (non-verbose).

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createPassFailTree(t, 1, 0)
    req.Args = []string{"test", testDir}
    return nil
}
```
