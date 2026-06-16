# Scenario

**Feature**: a temporary test tree with 2 passing and 1 failing leaves is created

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A temporary test tree with 2 passing and 1 failing leaves is created.

## Steps
1. Create test tree with 2 pass + 1 fail.
2. Run `doctest test <dir>` (non-verbose).

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createPassFailTree(t, 2, 1)
    req.Args = []string{"test", testDir}
    return nil
}
```
