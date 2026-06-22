# Helper Shadow Fixture

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system** — under test.

### Behaviors
- **run** — executes the scenario.

```go
import "testing"

type Request struct {
	Name string
}

type Response struct {
	Message string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	return &Response{Message: "hello " + req.Name}, nil
}
```