# Scenario

**Feature**: Inline any-of on a single line

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Templates use inline `<any-of>`.
