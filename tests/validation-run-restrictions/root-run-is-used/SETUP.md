# Scenario

**Feature**: rule R4: Only the root's Run is used in generated code

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root may omit Setup (prose-only OK)
```

## Preconditions
- Rule R4: Only the root's Run is used in generated code.
- The root Run function is the one executed, not any overrides that may exist.
