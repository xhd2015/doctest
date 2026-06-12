## Setup

- Defines a helper function `myHelper` at the root level.
- Defines the shared Request and Response model.
- A child SETUP.md also defines `myHelper`, causing a name collision
  when both are emitted as closures with `:=`.

```go
type Request struct {
    Name string
}

type Response struct {
    Message string
}

func myHelper(s string) string {
    return "root: " + s
}

func Setup(t *testing.T, req *Request) error {
    req.Name = "root"
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    return &Response{Message: "hello " + req.Name}, nil
}
```
