# Deeply Nested Root — DOCTEST.md Boundary

## Version
0.0.2


This is a deeply nested self-contained test root. It verifies that
DOCTEST.md boundaries work at any depth — even inside another nested root.

This root defines its own `Request{ID, Data}` and `Response{Status, Message}`
types, which are completely different from any ancestor types.

## How to Run

```sh
doctest test ./tests/boundary/nested/error_edge/
```

```go
import (
	"fmt"
	"testing"
)

type Request struct {
	ID	int
	Data	string
}
type Response struct {
	Status	string
	Message	string
}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	if req.ID <= 0 {
		return nil, fmt.Errorf("ID must be positive")
	}
	if req.Data == "" {
		return nil, fmt.Errorf("Data is required")
	}
	return &Response{Status: "ok", Message: "processed: " + req.Data}, nil
}
```
