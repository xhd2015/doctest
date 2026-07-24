# Root Run Used Fixture

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **fixture** — shared greet harness.

### Behaviors
- **run** — greets by name.

```go
import "fmt"
import "testing"

type Request struct {
	Name string
}

type Response struct {
	Message string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	return &Response{Message: fmt.Sprintf("hello %s", req.Name)}, nil
}
```