## Preconditions
- Parent git repo has go.mod with module `testproj`, no doctests.
- Inside parent, `sub/` is a separate git repo with go.mod `testproj/sub` and a DOCTEST.md tree.
- The nested module's path IS a child of the parent module path (so module path prefix matches).

## Steps
1. Create parent git repo with go.mod (module: testproj).
2. Create `sub/` — separate git repo with go.mod (module: testproj/sub) and a doctest tree.
3. Run `doctest test -v ./...` from parent.
4. Verify child's tests are NOT discovered (walk down stops at git boundary).

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    parent := t.TempDir()
    if err := initGitRepo(parent); err != nil {
        t.Fatalf("init parent repo: %v", err)
    }
    if err := writeGoMod(parent, "testproj"); err != nil {
        t.Fatalf("write parent go.mod: %v", err)
    }

    subDir := filepath.Join(parent, "sub")
    if err := os.MkdirAll(subDir, 0755); err != nil {
        t.Fatalf("mkdir sub: %v", err)
    }
    if err := initGitRepo(subDir); err != nil {
        t.Fatalf("init sub repo: %v", err)
    }
    if err := writeGoMod(subDir, "testproj/sub"); err != nil {
        t.Fatalf("write sub go.mod: %v", err)
    }
    if err := createTestTree(subDir, "child_test"); err != nil {
        t.Fatalf("create child_test: %v", err)
    }

    req.WorkDir = parent
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
