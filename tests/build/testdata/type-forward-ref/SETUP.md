## Setup

- Defines types in file-order: LocationEntry references GitInfo before GitInfo is declared.
- When emitted inside a function body, this causes "undefined: GitInfo".

```go
type LocationEntry struct {
    Path string
    Git  *GitInfo
}

type GitInfo struct {
    Type string
}

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
