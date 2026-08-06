# Scenario

**Feature**: bare `!` and postfix bang are parse errors

```
EvalLabelExpr("!", labels) -> error
EvalLabelExpr("e2e !", labels) -> error
```

## Steps

1. Expect non-nil parse errors.
