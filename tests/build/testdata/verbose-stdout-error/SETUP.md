# Scenario

**Feature**: tests for verbose stdout error

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Setup

- Minimal valid doc-style tree to trigger the `doctest test -v` code path.

```go
type Request struct {
    Name string
}

type Response struct {
    Message string
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Name = "default"
    return nil
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
    return &Response{Message: "hello " + req.Name}, nil
}
```
