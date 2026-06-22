# Scenario

**Feature**: the root DOCTEST.md defines Request, Response, and Run

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Child No Setup Fixture

## Preconditions
- The root DOCTEST.md defines Request, Response, and Run.
- The leaf SETUP.md defines neither Setup nor Run, which triggers a validation
  error because non-root SETUP.md must have func Setup.