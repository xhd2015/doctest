# Scenario

**Feature**: Contains order-free fragments

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Templates use `<contains>`.
