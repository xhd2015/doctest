# Scenario

**Feature**: tests for pkg1

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

```go
import "testing"

type Request struct {
    Name string
}

type Response struct {
    Name string
}

func Setup(t *testing.T, req *Request) error {
    req.Name = "pkg1"
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    return &Response{Name: req.Name}, nil
}
```
