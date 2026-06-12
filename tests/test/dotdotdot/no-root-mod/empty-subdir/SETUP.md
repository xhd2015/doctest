## Preconditions
- testdata/ is an empty directory with no go.mod and no DOCTest trees.

## Steps
1. Set WorkDir to a temp directory (empty, no go.mod anywhere).
2. Run `doctest test -v ./...`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.WorkDir = t.TempDir()
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
