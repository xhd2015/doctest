## Setup

- Minimal valid doc-style tree to trigger the `doctest test -v` code path.

```go
type Request struct {
    Name string
}

type Response struct {
    Message string
}

func Setup(t *testing.T, req *Request) error {
    req.Name = "default"
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    return &Response{Message: "hello " + req.Name}, nil
}
```
