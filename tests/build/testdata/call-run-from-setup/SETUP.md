# Scenario

**Feature**: tests for call run from setup

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Setup

- Defines the shared Request and Response model.
- Defines func Run which the child Setup will call.
- The generator lowers func Run to lowercase `run` closure,
  but the child Setup body still references uppercase `Run` → "undefined: Run".

```go
type Request struct {
    Name string
}

type Response struct {
    Message string
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Name = "root"
    return nil
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
    return &Response{Message: "hello " + req.Name}, nil
}
```
