# DOCTEST.md Boundary Handling Tests

## Version
0.0.2


These tests verify that `DOCTEST.md` creates an inheritance firewall:
no types, helpers, or setup functions cross the boundary. Each tree
rooted at a `DOCTEST.md` is a self-contained test tree with its own
`Request`/`Response` types.

## Decision Tree

```
tests/boundary/                         [parent root: Request{}, Response{}]
├── leaf/                               → uses parent types, Run stub → error
└── nested/                             [DOCTEST.md boundary]
    ├── DOCTEST.md                      → marks nested root
    ├── SETUP.md                        → defines Request{Name}, Response{Greeting}
    ├── happy/                          → uses nested types, Run succeeds
    └── error_edge/                     [DOCTEST.md boundary]
        ├── DOCTEST.md                  → marks deeply nested root
        ├── SETUP.md                    → defines Request{ID,Data}, Response{Status,Message}
        └── check/                      → uses deeply nested types, Run succeeds
```

## Test Index

| Leaf | Description |
|------|-------------|
| `leaf` | Parent root leaf: stub Run returns error, asserts error is non-nil |
| `nested/happy` | Nested root leaf: Run returns greeting for given Name |
| `nested/error_edge/check` | Deeply nested root leaf: Run processes ID and Data successfully |

## How to Run

```sh
doctest test ./tests/boundary/                       # parent root
doctest test ./tests/boundary/nested/                # nested root
doctest test ./tests/boundary/nested/error_edge/     # deeply nested root
doctest vet ./tests/boundary/                        # vet must skip nested roots
```

```go
import (
	"fmt"
	"testing"
)

type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) {
	return nil, fmt.Errorf("stub: not implemented")
}
```
