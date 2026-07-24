# Scenario

**Feature**: OR matches when any operand label is present

```
EvalLabelExpr("slow || heavy", {slow}) -> true
EvalLabelExpr("slow || heavy", {heavy}) -> true
EvalLabelExpr("slow || heavy", {fast}) -> false
```

## Steps

1. Assert OR semantics.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```