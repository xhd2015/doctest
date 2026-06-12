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

func Setup(t *testing.T, req *Request) error {
    req.Name = "root"
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    return &Response{Message: "hello " + req.Name}, nil
}
```
