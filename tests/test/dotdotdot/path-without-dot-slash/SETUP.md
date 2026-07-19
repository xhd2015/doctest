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
1. Set WorkDir to the testdata/ directory.
2. Run `doctest test -v alpha/...` (without `./` prefix).

```go
import (
"github.com/xhd2015/doctest/session"
"path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.WorkDir = filepath.Join(d.DOCTEST_CASE, "testdata")
    req.Args = []string{"test", "-v", "alpha/..."}
    return nil
}
```
