# Scenario

**Feature**: a temporary project with go.mod + multiple DOCTEST.md trees exists

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A temporary project with go.mod + multiple DOCTEST.md trees exists.

## Steps
1. Create temp project via parent helper.
2. Run `doctest test ./...` from the project root.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createTempProject(t, req)
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
