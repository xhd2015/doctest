# Scenario

**Feature**: a project with complex doctest tree structure (ancestor/ leaf/ nested sub2/)

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A project with complex doctest tree structure (ancestor/ leaf/ nested sub2/).

## Steps
1. Create the mixed test project.
2. Run `doctest test -v ./ancestor/leaf/...` from the project root.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createMixedTestProject(t)
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./ancestor/leaf/..."}
    return nil
}
```
