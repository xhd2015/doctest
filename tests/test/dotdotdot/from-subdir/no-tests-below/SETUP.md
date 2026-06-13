## Preconditions
- A temporary project with go.mod but no DOCTEST.md trees at or below the working directory.
- A DOCTEST.md tree exists elsewhere in the module (above the working directory).

## Steps
1. Create temp project with doctest trees only above the working directory.
2. Run `doctest test ./...` from a subdirectory with no DOCTEST.md or its subdirectories.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createTempProject(t, req)

    // run from alpha_test/simple/ — no DOCTEST.md at or below this dir
    req.WorkDir = filepath.Join(projDir, "alpha_test", "simple")
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
