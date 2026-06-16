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
1. Set WorkDir to `../basic/testdata` (relative from this leaf).
2. Run `doctest test -v ./alpha/simple/...`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.WorkDir = "../basic/testdata"
    req.Args = []string{"test", "-v", "./alpha/simple/..."}
    return nil
}
```
