# Scenario

**Feature**: `--changed` in subcommand help via in-process CLI (no product binary)

```
cli.RunWithWriter -> doctest <subcmd> --help -> stdout lists --changed
```

## Preconditions

- Help is covered in-process via `cli.RunWithWriter` (same usage strings as the product binary).
- Unlabeled (fast); no `testbin`, no `label: e2e`.
- Policy selection stays in-process under `git-context/in-git-repo/`.

## Steps

1. Root help Setup is a no-op (no binary).
2. Leaf sets `Args` to `<subcmd> --help`.
3. `Run` calls `cli.RunWithWriter` when Args are set without TreeDir.

## Context

- No fixture tree required for help.
