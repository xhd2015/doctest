# Scenario

**Feature**: malformed expressions return parse errors

```
EvalLabelExpr("slow &&", labels) -> error
```

## Steps

1. Expect non-nil error from evaluator.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```