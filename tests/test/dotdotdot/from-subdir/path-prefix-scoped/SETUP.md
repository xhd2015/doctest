## Preconditions
- A temporary project with `parent_subpath/` doctest tree at module root.
- A `workdir/` subdirectory that also contains a `local_subpath/` doctest tree.
- Running `./local_subpath/...` from `workdir/` must only find the CWD-relative one.

## Steps
1. Create temp project with doctest trees at both module-root and CWD-relative levels.
2. Run `doctest test ./local_subpath/...` from `workdir/`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createTempProject(t, req)

    // doctest tree at module root level (should NOT be found)
    if err := createTestTree(projDir, "parent_subpath"); err != nil {
        t.Fatalf("create parent_subpath: %v", err)
    }

    // subdirectory with its own doctest tree (the one that should be found)
    workDir := filepath.Join(projDir, "workdir")
    if err := createTestTree(workDir, "local_subpath"); err != nil {
        t.Fatalf("create local_subpath: %v", err)
    }

    req.WorkDir = workDir
    req.Args = []string{"test", "-v", "./local_subpath/..."}
    return nil
}
```
