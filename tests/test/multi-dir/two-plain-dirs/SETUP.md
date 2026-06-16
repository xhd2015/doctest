# Scenario

**Feature**: a temp project exists with `test_a` and `test_b` doctest trees

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A temp project exists with `test_a` and `test_b` doctest trees.

## Steps
1. Create the temp project.
2. Run `doctest test test_a test_b` (two plain directory arguments).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    projDir := createMultiDirProject(t, req)
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "test_a", "test_b"}
    return nil
}
```
