# Scenario

**Feature**: short-path in-process CLI for `doctest metrics` help / unknown

```
cli.RunWithWriter -> doctest metrics --help | unknown subcommand
  -> usage text / non-zero exit (no product binary, no testbin)
```

## Preconditions

- Help/unknown use `cli.RunWithWriter` (Parallel-safe; no `os.Setenv`/`Chdir`).
- Unlabeled (fast). Analyze combinatorics stay under path/last/top/….

## Steps

1. Leaf sets `Args` for help or unknown subcommand.
2. Root `Run` detects short CLI args and dispatches via `cli.RunWithWriter`.

## Context

- No MetricsRoot fixtures required for help/unknown.
