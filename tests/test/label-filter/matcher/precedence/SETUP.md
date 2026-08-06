# Scenario

**Feature**: AND binds tighter than OR

```
a || b && c  ≡  a || (b && c)
```

## Steps

1. Compare match results for equivalent label sets.
