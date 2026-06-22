# Child Setup Only Valid Fixture

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **fixture** — multi-level setup chain.

### Behaviors
- **run** — scales value by 10.

```go
import "fmt"
import "testing"

type Request struct {
	Value int
}

type Response struct {
	Result int
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Value < 0 {
		return nil, fmt.Errorf("Value must be non-negative")
	}
	return &Response{Result: req.Value * 10}, nil
}
```