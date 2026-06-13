## Preconditions
- A temporary project with go.mod + two DOCTEST.md trees (alpha_test, beta_test).
- A nested submodule with its own DOCTEST.md tree (hidden_test).

## Steps
1. Create temp project via `createTempProject`.
2. Run `doctest test ./...` from `alpha_test/` (which has its own DOCTEST.md).

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createTempProject(t, req)
    req.WorkDir = filepath.Join(projDir, "alpha_test")
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
