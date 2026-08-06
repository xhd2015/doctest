# Scenario

**Feature**: surrounding whitespace is ignored when parsing

```
" slow && ui " matches {slow,ui}
```

## Steps

1. Evaluate trimmed expression string.
