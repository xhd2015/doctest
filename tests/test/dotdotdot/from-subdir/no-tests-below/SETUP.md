# Scenario

**Feature**: a temporary project with go.mod but no DOCTEST.md trees at or below the working directory

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A temporary project with go.mod but no DOCTEST.md trees at or below the working directory.
- A DOCTEST.md tree exists elsewhere in the module (above the working directory).

## Steps
1. Create temp project with doctest trees only above the working directory.
2. Run `doctest test ./...` from a subdirectory with no DOCTEST.md or its subdirectories.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    projDir := createTempProject(t, req)

    // run from alpha_test/simple/ — no DOCTEST.md at or below this dir
    req.WorkDir = filepath.Join(projDir, "alpha_test", "simple")
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
