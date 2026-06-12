## Preconditions
- A minimal valid doctest tree is available.

## Steps
1. Run `doctest validate <dir>`.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    dir := t.TempDir()
    if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# tests\n"), 0644); err != nil {
        t.Fatal(err)
    }
    if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte("setup\n"), 0644); err != nil {
        t.Fatal(err)
    }
    req.Args = []string{"validate", dir}
    return nil
}
```
