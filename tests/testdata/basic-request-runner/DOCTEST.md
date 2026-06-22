# Basic Request Runner Fixture

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **runner** — shared Request/Response harness.

### Behaviors
- **run** — dispatches greet/fail actions.

```go
import "fmt"
import "testing"

type Request struct {
	Action string
	Name   string
}

type Response struct {
	Message string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch req.Action {
	case "greet":
		return &Response{Message: fmt.Sprintf("hello %s", req.Name)}, nil
	case "fail":
		return nil, fmt.Errorf("requested failure for %s", req.Name)
	default:
		return nil, fmt.Errorf("unknown action %q", req.Action)
	}
}
```