# Scenario

**Feature**: testdata/ has no go.mod but contains two DOCTest trees: tests/my-feature/ and other/other-feature/. Running tests/... should only find tests/my-feature/.

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
- tests/my-feature/ and other/other-feature/ are valid DOCTest trees.

## Steps
1. Copy testdata/ to a temp directory (outside any Go module).
2. Run `doctest test -v tests/...` from the temp directory.

```go
import (
    "os/exec"
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
    // init git in case the working directory is inside a git repo
    exec.Command("git", "-C", tmpTestData, "init").CombinedOutput()
    req.WorkDir = tmpTestData
    req.Args = []string{"test", "-v", "tests/..."}
    return nil
}
```
