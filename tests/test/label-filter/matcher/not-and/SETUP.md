# Scenario

**Feature**: `!` binds tighter than `&&`

```
EvalLabelExpr("!e2e && flaky", {flaky}) -> true
EvalLabelExpr("!e2e && flaky", {e2e, flaky}) -> false
EvalLabelExpr("!e2e && flaky", {}) -> false
```

## Steps

1. Assert combined negation and AND.
