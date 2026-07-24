# Child No Setup Fixture

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **fixture** — shared harness.

### Behaviors
- **run** — doubles non-negative values.

```go
import "fmt"
import "testing"

type Request struct {
	Value int
}

type Response struct {
	Result int
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	if req.Value < 0 {
		return nil, fmt.Errorf("Value must be non-negative")
	}
	return &Response{Result: req.Value * 2}, nil
}
```