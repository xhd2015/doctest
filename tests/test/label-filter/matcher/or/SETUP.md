# Scenario

**Feature**: OR matches when any operand label is present

```
EvalLabelExpr("slow || flaky", {slow}) -> true
EvalLabelExpr("slow || flaky", {flaky}) -> true
EvalLabelExpr("slow || flaky", {fast}) -> false
```

## Steps

1. Assert OR semantics.
