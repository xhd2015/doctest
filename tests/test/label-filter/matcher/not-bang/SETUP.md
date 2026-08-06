# Scenario

**Feature**: unary `!` negates a label atom (includes unlabeled)

```
EvalLabelExpr("!e2e", {}) -> true
EvalLabelExpr("!e2e", {heavy}) -> true
EvalLabelExpr("!e2e", {e2e}) -> false
EvalLabelExpr("!e2e", {e2e, heavy}) -> false
```

## Steps

1. Assert bang-negation semantics on empty and non-empty label sets.
