# Scenario

**Feature**: rule R3: Non-root SETUP.md files cannot redefine func Run

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root may omit Setup (prose-only OK)
```

## Preconditions
- Rule R3: Non-root SETUP.md files cannot redefine func Run.
- Run is reserved for the root SETUP.md only.
