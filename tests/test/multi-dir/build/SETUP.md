# Scenario

**Feature**: build tests have a longer timeout due to compilation

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- Build tests have a longer timeout due to compilation.

## Steps
1. Extend timeout for build operations.

```go
import (
    "testing"
    "time"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
