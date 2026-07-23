# Scenario

**Feature**: `!` binds tighter than `&&`

```
EvalLabelExpr("!e2e && heavy", {heavy}) -> true
EvalLabelExpr("!e2e && heavy", {e2e, heavy}) -> false
EvalLabelExpr("!e2e && heavy", {}) -> false
```

## Steps

1. Assert combined negation and AND.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```
