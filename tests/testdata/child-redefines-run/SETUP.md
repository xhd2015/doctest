# Scenario

**Feature**: the root DOCTEST.md defines Request, Response, and Run

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Child Redefines Run Fixture

## Preconditions
- The root DOCTEST.md defines Request, Response, and Run.
- A child SETUP.md also defines Run, which should trigger a validation error.