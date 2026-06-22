# Nested Root — DOCTEST.md Boundary

## Version
0.0.2


This is a self-contained test root nested inside `tests/boundary/`.
The parent root must skip this directory entirely due to the
`DOCTEST.md` inheritance firewall.

This root defines its own `Request{Name}` and `Response{Greeting}`
types, which differ from the parent's empty `Request{}`/`Response{}`.

## Decision Tree

```
nested/                         [self-contained root: Request{Name}, Response{Greeting}]
└── happy/                      → sets Name, asserts Greeting is non-empty
```

## How to Run

```sh
doctest test ./tests/boundary/nested/
```

```go
import (
	"fmt"
	"testing"
)

type Request struct {
	Name string
}
type Response struct {
	Greeting string
}
func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("Name is required")
	}
	return &Response{Greeting: "Hello, " + req.Name + "!"}, nil
}
```
