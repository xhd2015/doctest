# Scenario

**Feature**: rule R2: Non-root SETUP.md may omit func Setup (prose-only org nodes)

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root missing Setup -> build error
```

## Preconditions
- Rule R2: Non-root SETUP.md may omit func Setup (prose-only org nodes).
- Run is reserved for root; non-root SETUP.md without at least Setup is invalid.
