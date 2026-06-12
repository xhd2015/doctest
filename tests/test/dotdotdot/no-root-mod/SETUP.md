## Group: no-root-mod
Tests for `./...` when CWD has no `go.mod`.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func copyDir(dst, src string) error {
    return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        rel, _ := filepath.Rel(src, path)
        if rel == "." {
            return nil
        }
        target := filepath.Join(dst, rel)
        if d.IsDir() {
            return os.MkdirAll(target, 0755)
        }
        if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
            return err
        }
        data, readErr := os.ReadFile(path)
        if readErr != nil {
            return readErr
        }
        return os.WriteFile(target, data, 0644)
    })
}

func Setup(t *testing.T, req *Request) error {
    t.Logf("no-root-mod group: WorkDir=%s", req.WorkDir)
    return nil
}
```
