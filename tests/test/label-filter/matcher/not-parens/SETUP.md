# Scenario

**Feature**: `!(…)` groups OR under negation

```
EvalLabelExpr("!(e2e || flaky)", {}) -> true
EvalLabelExpr("!(e2e || flaky)", {slow}) -> true
EvalLabelExpr("!(e2e || flaky)", {e2e}) -> false
EvalLabelExpr("!(e2e || flaky)", {flaky}) -> false
```

## Steps

1. Assert parenthesized negation.
