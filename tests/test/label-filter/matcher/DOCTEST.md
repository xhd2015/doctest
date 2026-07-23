# Label Expression Matcher

## Version

0.0.2

**Layer L2 in-process** — nested library tree for `core.EvalLabelExpr`.
Leaves call the evaluator in-process via Assert helpers. **No product binary**,
**no `label: heavy`**. Discovered and run by default (`doctest test` without
`--label`).

# DSN (Domain Specific Notion)

The **label expression evaluator** (`core.EvalLabelExpr`) parses a boolean
expression over leaf label sets and returns whether a label set satisfies it.

Nested tree: library contract only (no subprocess).

## Decision Tree

```
matcher/                              [L2 in-process — EvalLabelExpr]
├── single-label/                     atom match / miss
├── and/                              && requires all tokens
├── or/                               || any token
├── precedence/                       && binds tighter than ||
├── parentheses/                      grouping overrides
├── whitespace/                       trim around expression
├── invalid-syntax/                   trailing && → parse error
├── not-bang/                         !e2e (includes unlabeled)
├── not-and/                          !e2e && heavy
├── not-parens/                       !(e2e || flaky)
└── not-invalid/                      bare ! / trailing ! / "not e2e"
```


## How to Run

```sh
# default discovery — all matcher leaves (unlabeled L2)
doctest test ./tests/test/label-filter/matcher
doctest test ./tests/test/label-filter/matcher/...
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
