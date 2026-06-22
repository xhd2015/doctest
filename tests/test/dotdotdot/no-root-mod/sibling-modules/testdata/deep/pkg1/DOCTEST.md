# pkg1 Tests

## Version
0.0.2


```go
import "testing"

type Request struct {
	Name string
}
type Response struct {
	Name string
}
func Run(t *testing.T, req *Request) (*Response, error) {
	return &Response{Name: req.Name}, nil
}
```
