# Scenario

**Feature**: `doctest vet` rejects vacuous `func Setup` bodies and allows prose-only SETUP.md (no Go block)

```
# vacuous Setup (non-root) hard-fails
leaf SETUP with only return nil | blank assigns + return nil
  -> runner.VetArgs -> non-zero
  -> message: remove Go code block (not "implement the behavior")

# prose-only SETUP (no go fence) allowed
intermediate or leaf SETUP with # Scenario prose only
  -> exit 0

# real Setup with real work allowed
Setup sets a field / does work -> exit 0
```

## Preconditions

- Fixtures are minimal DOCTEST trees under `t.TempDir()` (not the repo tree).
- Vacuous / real Setup fixtures put the target `SETUP.md` under a **non-root**
  path (`leaf/` or `group/`) so `CheckSetupBodyNotStub` applies (root SETUP is
  not stub-checked today).
- Prose-only fixtures omit the go code fence entirely.
- In-process via shared root `Run` → `runner.VetArgs` (no product binary).
- Parallel-safe: no Setenv/Chdir/stdio reassignment.
- Grouping Setup is a no-op (leaves own fixture work). After product relaxes
  prose-only discovery, this intermediate SETUP may become prose-only.

## Steps

1. Grouping is organization-only (prose SETUP, no Go block); each leaf builds its fixture and sets Args.
2. Run via shared root `Run` → in-process `VetArgs`.
