# Scenario

**Feature**: the root Setup has built the doctest binary and set a generous timeout

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The root Setup has built the doctest binary and set a generous timeout.
- Child leaves test the `--gen-dir` flag behavior.

## Steps
1. Ensure the timeout from the parent Setup is preserved (no override needed at this grouping level).

```go
import (
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
