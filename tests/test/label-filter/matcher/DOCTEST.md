# Label Expression Matcher

## Version

0.0.2

# DSN (Domain Specific Notion)

The **label expression evaluator** (`core.EvalLabelExpr`) parses a boolean
expression over leaf label sets and returns whether a label set satisfies it.

Nested tree: library contract only (no subprocess).

## How to Run

```sh
doctest test ./tests/test/label-filter/matcher
```

```go
import (
	"testing"
)

type Request struct{}

type Response struct{}

func Run(t *testing.T, req *Request) (*Response, error) {
	_ = req
	return &Response{}, nil
}
```