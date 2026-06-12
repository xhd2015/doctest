## Preconditions
- testdata/ contains a directory with .gitignore and multiple DOCTest trees (one matches .gitignore).
- There is NO git repository (no `git init`).

## Steps
1. Copy testdata/ to a temp directory outside any Go module.
2. Do NOT run `git init` — this tests that .gitignore is NOT respected without a git repo.
3. Run `doctest test -v ./...`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    srcTestData := "./testdata"
    tmpTestData := filepath.Join(t.TempDir(), "testdata")
    if err := copyDir(tmpTestData, srcTestData); err != nil {
        t.Fatalf("copy testdata: %v", err)
    }
    req.WorkDir = tmpTestData
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
