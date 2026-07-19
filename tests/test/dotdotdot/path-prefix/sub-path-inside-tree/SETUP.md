# Scenario

**Feature**: the `basic/testdata/` has go.mod with DOCTest trees: alpha/ and beta/

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The `basic/testdata/` has go.mod with DOCTest trees: alpha/ and beta/.
- The base path `alpha/simple/` is a subdirectory within the `alpha/` doctest tree (no DOCTEST.md at that level).

## Steps
1. Set WorkDir to the abs path of sibling `basic/testdata` (via `d.DOCTEST_CASE`; cwd is undetermined).
2. Run `doctest test -v ./alpha/simple/...`.

```go
import (
    "path/filepath"
    "testing"

    "github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    // Sibling of this leaf: path-prefix/basic/testdata
    req.WorkDir = filepath.Join(d.DOCTEST_CASE, "..", "basic", "testdata")
    req.Args = []string{"test", "-v", "./alpha/simple/..."}
    return nil
}
```
