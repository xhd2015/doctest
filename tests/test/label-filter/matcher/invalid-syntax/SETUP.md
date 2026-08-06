# Scenario

**Feature**: malformed expressions return parse errors

```
EvalLabelExpr("slow &&", labels) -> error
```

## Steps

1. Expect non-nil error from evaluator.
