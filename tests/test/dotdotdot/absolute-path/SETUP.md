# Scenario

**Feature**: absolute path `/...` pattern support

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Group: absolute-path
Tests for absolute `<prefix>/...` pattern support (same semantics as `./<prefix>/...`).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    t.Logf("absolute-path group")
    return nil
}
```