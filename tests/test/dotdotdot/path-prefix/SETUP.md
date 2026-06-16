# Scenario

**Feature**: tests for path prefix

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Group: path-prefix
Tests for `./<prefix>/...` pattern support.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    t.Logf("path-prefix group: WorkDir=%s", req.WorkDir)
    return nil
}
```
