# Scenario

**Feature**: rule R2: Non-root SETUP.md may omit func Setup (prose-only org nodes)

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root may omit Setup (prose-only OK)
```

## Preconditions
- Rule R2: Non-root SETUP.md may omit func Setup (prose-only org nodes).
- Run is reserved for root; non-root SETUP may be prose-only without Setup.
