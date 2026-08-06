# Scenario

**Feature**: help surfaces for `cache` and top-level command list

```
# help path (no cache filesystem)
cli.RunWithWriters -> doctest cache --help | doctest --help
  -> usage / command list on stdout
```

## Preconditions

- No CacheHome isolation required for help-only leaves.
- L2 in-process CLI; unlabeled.

## Steps

1. Each leaf sets `req.Args` for the help variant under test.

## Context

- Grouping only; leaves hold Assert.
