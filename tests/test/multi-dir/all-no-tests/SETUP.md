# Scenario

**Feature**: a temp project exists with only a non-doctest directory (`no_tests`). No valid test trees are targeted

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A temp project exists with only a non-doctest directory (`no_tests`). No valid test trees are targeted.

## Steps
1. Create the temp project.
2. Run `doctest test no_tests` — the dir has no doctest tree, so `ErrNoTestsFound` is returned.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    projDir := createMultiDirProject(t, req)
    req.WorkDir = projDir
    req.Args = []string{"test", "no_tests"}
    return nil
}
```
