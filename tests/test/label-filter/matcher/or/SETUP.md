# Scenario

**Feature**: OR matches when any operand label is present

```
EvalLabelExpr("slow || heavy", {slow}) -> true
EvalLabelExpr("slow || heavy", {heavy}) -> true
EvalLabelExpr("slow || heavy", {fast}) -> false
```

## Steps

1. Assert OR semantics.
