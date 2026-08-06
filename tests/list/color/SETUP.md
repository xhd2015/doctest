# Scenario

**Feature**: color flags for list output (flag-only; no process Setenv)

```
Harness -> list [--color|--no-color] <root>
  -> ANSI gray on meta when forced on; none when default pipe or --no-color
  -> conflict of both flags -> error
```

## Preconditions

- Capture writers are non-TTY buffers → default auto color is off.
- Prefer `--color` / `--no-color` over `NO_COLOR` env for parallel safety.

## Steps

1. Grouping Setup is a no-op.
2. Leaves share a one-leaf fixture pattern via their own Setup.
