## Preconditions
- testdata/ contains a directory with .gitignore, multiple subdirs (some gitignored), each with a DOCTest tree and go.mod.

## Steps
1. Copy testdata/ to a temp directory outside any Go module.
2. Run `git init` in the temp directory so .gitignore is respected.
3. Run `doctest test -v ./...`.

```go
import (
    "os/exec"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    srcTestData := "./testdata"
    tmpTestData := filepath.Join(t.TempDir(), "testdata")
    if err := copyDir(tmpTestData, srcTestData); err != nil {
        t.Fatalf("copy testdata: %v", err)
    }
    if out, err := exec.Command("git", "-C", tmpTestData, "init").CombinedOutput(); err != nil {
        t.Fatalf("git init: %v\n%s", err, out)
    }
    req.WorkDir = tmpTestData
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
