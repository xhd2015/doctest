# Scenario

**Feature**: invalid `-timeout` values are rejected at parse time

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- An invalid duration string is passed to `-timeout`.

## Steps
1. Run `doctest test -timeout bogus .` (directory unused — parse fails first).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"test", "-timeout", "bogus", "."}
    return nil
}
```
