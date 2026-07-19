# Scenario

**Feature**: testdata/ contains two modules (mod_a, mod_b) with DOCTest trees

## Preconditions
- testdata/ contains two modules (mod_a, mod_b) with DOCTest trees.
- There is NO git repository (no `git init`).

## Steps
1. Copy testdata/ to a temp directory outside any Go module.
2. Do NOT run `git init` — this tests that no git fatal message appears when not in a git repo.
3. Run `doctest test -v ./...`.

```go
import (
"github.com/xhd2015/doctest/session"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    srcTestData := filepath.Join(d.DOCTEST_CASE, "testdata")
    tmpTestData := filepath.Join(t.TempDir(), "testdata")
    if err := copyDir(tmpTestData, srcTestData); err != nil {
        t.Fatalf("copy testdata: %v", err)
    }
    req.WorkDir = tmpTestData
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
