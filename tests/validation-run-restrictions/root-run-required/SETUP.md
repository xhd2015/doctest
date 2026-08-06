# Scenario

**Feature**: rule R1: Root SETUP.md must define type Request, type Response, and func Run

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root may omit Setup (prose-only OK)
```

## Preconditions
- Rule R1: Root SETUP.md must define type Request, type Response, and func Run.
- Without Run at the root, tree discovery fails.
