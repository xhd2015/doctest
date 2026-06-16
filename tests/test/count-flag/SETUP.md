# Scenario

**Feature**: a valid doc-style test tree exists in the repository

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A valid doc-style test tree exists in the repository.

## Steps
1. Run `doctest test <dir> -count=1 -v`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
    req.Args = []string{"test", exampleDir, "-count=1", "-v"}
    return nil
}
```
