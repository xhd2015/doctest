# Scenario

**Feature**: per-root leaf counts, L2:L3 split, and label distribution

```
Harness -> write root with labeled/unlabeled ASSERT leaves
  -> list <root>
  -> L2:L3=a:b (p2%/p3%) + labelDist (incl unlabeled)
```

## Preconditions

- Each leaf builds a single-root fixture tailored to label/L2:L3 semantics.
- L3 identity is label `e2e` only; cost labels like `heavy` without `e2e` stay L2.

## Steps

1. Grouping Setup is a no-op.
2. Leaves write fixtures and set Args.
