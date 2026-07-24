# Scenario

**Feature**: `!(…)` groups OR under negation

```
EvalLabelExpr("!(e2e || flaky)", {heavy}) -> true
EvalLabelExpr("!(e2e || flaky)", {e2e}) -> false
EvalLabelExpr("!(e2e || flaky)", {flaky}) -> false
```

## Steps

1. Assert parenthesized negation.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```
