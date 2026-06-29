# Scenario

**Feature**: testdata/ has go.mod at root with two DOCTest trees: alpha/ and beta/

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- testdata/ has go.mod at root with two DOCTest trees: alpha/ and beta/.

## Steps
1. Resolve testdata/ to an absolute path.
2. Run `doctest test -v <abs-testdata>/alpha/...` (shell-style `/...` suffix on absolute path).

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    testdata, err := filepath.Abs(filepath.Join(DOCTEST_ROOT, "test", "dotdotdot", "path-prefix", "basic", "testdata"))
    if err != nil {
        return err
    }
    req.WorkDir = testdata
    alphaPattern := filepath.ToSlash(filepath.Join(testdata, "alpha")) + "/..."
    req.Args = []string{"test", "-v", alphaPattern}
    return nil
}
```