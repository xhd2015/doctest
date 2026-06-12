## Preconditions
- The target path exists but is not a directory.

## Steps
1. Run `doctest agent fill-code <file>`.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    path := filepath.Join(t.TempDir(), "target.txt")
    if err := os.WriteFile(path, []byte("not a dir"), 0644); err != nil {
        t.Fatal(err)
    }
    req.Args = []string{"agent", "fill-code", path}
    return nil
}
```
