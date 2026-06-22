# Child Redefines Run Fixture

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **fixture** — shared types with child Run violation.

### Behaviors
- **run** — root Run implementation.

```go
import "fmt"
import "testing"

type Request struct {
	Action string
}

type Response struct {
	Result string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Action == "" {
		return nil, fmt.Errorf("Action is required")
	}
	return &Response{Result: "root:" + req.Action}, nil
}
```