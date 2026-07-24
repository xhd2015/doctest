# Scenario

**Feature**: AND requires every operand label on the leaf

```
EvalLabelExpr("slow && ui", {slow,ui}) -> true
EvalLabelExpr("slow && ui", {slow}) -> false
```

## Steps

1. Assert AND semantics.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```