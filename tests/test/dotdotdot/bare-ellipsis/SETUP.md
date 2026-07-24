# Scenario

**Feature**: the parent dotdotdot helper functions are available

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The parent dotdotdot helper functions are available.

## Steps
1. Run `doctest test ...` (bare ellipsis, no `./` prefix, no path).
2. This should produce an error because bare `...` is not supported.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.WorkDir = t.TempDir()
    req.Args = []string{"test", "-v", "..."}
    return nil
}
```
