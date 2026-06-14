## Preconditions
- Same git repo with go.mod at root.
- Root has a doctest tree (parent_test).
- `child/` subdir has no doctests.

## Steps
1. Create git repo with go.mod and a doctest tree at root.
2. Run `doctest test -v ./...` from child/.
3. Verify no tests are found (`./...` only looks down from CWD, and child has no tests).

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    root := t.TempDir()
    if err := initGitRepo(root); err != nil {
        t.Fatalf("init repo: %v", err)
    }
    if err := writeGoMod(root, "testproj"); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }
    if err := createTestTree(root, "parent_test"); err != nil {
        t.Fatalf("create parent_test: %v", err)
    }

    childDir := filepath.Join(root, "child")
    if err := os.MkdirAll(childDir, 0755); err != nil {
        t.Fatalf("mkdir child: %v", err)
    }

    req.WorkDir = childDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
