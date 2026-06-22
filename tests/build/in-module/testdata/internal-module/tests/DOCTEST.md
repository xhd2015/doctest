# Tests

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system** — under test.

### Behaviors
- **run** — executes the scenario.


```go
import (
    "testing"
    "example.com/app/internal/greet"
)

type Request struct{}
type Response struct{ Message string }

func Run(t *testing.T, req *Request) (*Response, error) {
    return &Response{Message: greet.Hello()}, nil
}
```