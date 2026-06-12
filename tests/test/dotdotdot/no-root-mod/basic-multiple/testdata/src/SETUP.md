```go
import "testing"

type Request struct {
    Name string
}

type Response struct {
    Name string
}

func Setup(t *testing.T, req *Request) error {
    req.Name = "src"
    return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
    return &Response{Name: req.Name}, nil
}
```
