## Preconditions
- Same git repo with go.mod at root.
- Root has no doctests.
- `child/` subdir has a doctest tree.

## Steps
1. Create git repo with go.mod (no doctests at root).
2. Create a doctest tree in `child/` subdir.
3. Run `doctest test -v ./...` from repo root.
4. Verify child tests are discovered (walk down within same repo should find them).

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

    childDir := filepath.Join(root, "child")
    if err := createTestTree(childDir, "child_test"); err != nil {
        t.Fatalf("create child_test: %v", err)
    }

    req.WorkDir = root
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
