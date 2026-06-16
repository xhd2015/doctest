# Scenario

**Feature**: testdata/ is an empty directory with no go.mod and no DOCTest trees

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- testdata/ is an empty directory with no go.mod and no DOCTest trees.

## Steps
1. Set WorkDir to a temp directory (empty, no go.mod anywhere).
2. Run `doctest test -v ./...`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.WorkDir = t.TempDir()
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
