# Scenario

**Feature**: single-token expression matches one label

```
EvalLabelExpr("slow", {slow}) -> true
EvalLabelExpr("slow", {}) -> false
EvalLabelExpr("slow", {fast}) -> false
```

## Steps

1. Assert match outcomes for three label sets.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```