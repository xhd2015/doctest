# Scenario

**Feature**: parentheses override default precedence

```
(slow || flaky) && ui  matches only when ui plus (slow or flaky)
```

## Steps

1. Assert grouped OR then AND.
