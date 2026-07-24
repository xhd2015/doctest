# Scenario

**Feature**: the test command supports sub-directory resolution

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The test command supports sub-directory resolution.

## Steps
1. Create a temporary doc-style test tree.
2. Run `doctest test <sub-dir>` on the tree.

```go
import (
    "testing"
    "time"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Timeout = 20 * time.Second
    return nil
}
```
