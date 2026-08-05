# Scenario

**Feature**: the vet command detects anti-patterns in test file content (L2 in-process)

```
# inspect fixture tree for structural issues
runner.VetArgs([dir]) -> validate.RunWithOptions -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
| Parallel-unsafe: os.Setenv/Unsetenv/Clearenv | os.Chdir | t.Setenv | t.Chdir
| os.Stdout/Stderr/Stdin reassign | syscall.Setenv/Unsetenv
# allowed (positive control)
fmt.Fprint(os.Stdout, ...) without reassignment
```

## Preconditions

- Anti-pattern leaves write bad `SETUP.md` content (or place it under `testdata/`)
  and invoke in-process `VetArgs`.
- Parallel-unsafe leaves put the forbidden API only in **fixture** harness under
  `t.TempDir()` — the leaf harness that runs the scenario must not call
  Setenv/Chdir/stdio reassignment itself.

## Steps

1. Child leaf builds a fixture tree under `t.TempDir()`.
2. Run vet via shared in-process `Run` (no binary).

Organization-only grouping: no Go code block (leaves own fixture Setup).
