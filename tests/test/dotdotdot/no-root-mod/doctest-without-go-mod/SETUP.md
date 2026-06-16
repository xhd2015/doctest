# Scenario

**Feature**: testdata/ has no go.mod anywhere but contains a DOCTest tree under tests/

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- testdata/ has no go.mod anywhere.
- testdata/tests/my-feature/ is a valid DOCTest tree (DOCTEST.md + SETUP.md + simple/).

## Steps
1. Copy testdata/ to a temp directory (outside any Go module).
2. Run `doctest test -v ./...` from the temp directory.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    srcTestData := "./testdata"
    tmpDir := t.TempDir()
    tmpTestData := filepath.Join(tmpDir, "testdata")
    if err := copyDir(tmpTestData, srcTestData); err != nil {
        t.Fatalf("copy testdata: %v", err)
    }
    req.WorkDir = tmpTestData
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
