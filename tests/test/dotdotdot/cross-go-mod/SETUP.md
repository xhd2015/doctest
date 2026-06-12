## Preconditions
- The parent dotdotdot helpers (createTestTree, etc.) are available.

## Steps
1. Create a temp project with root go.mod and a nested go.mod whose module path is configurable.
2. Run `doctest test ./...` and verify which test trees are discovered.

```go
import (
    "os"
    "path/filepath"
    "testing"
    "time"
)

func createCrossModuleProject(t *testing.T, nestedSubDir, nestedModulePath, childTestName string) string {
    t.Helper()
    tmp := t.TempDir()

    if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write root go.mod: %v", err)
    }

    if err := createTestTree(tmp, "parent_test"); err != nil {
        t.Fatalf("create parent_test: %v", err)
    }

    nestedRoot := filepath.Join(tmp, nestedSubDir)
    if err := os.MkdirAll(nestedRoot, 0755); err != nil {
        t.Fatalf("mkdir nested: %v", err)
    }
    if err := os.WriteFile(filepath.Join(nestedRoot, "go.mod"), []byte("module "+nestedModulePath+"\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write nested go.mod: %v", err)
    }
    if err := createTestTree(nestedRoot, childTestName); err != nil {
        t.Fatalf("create %s: %v", childTestName, err)
    }

    return tmp
}

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
