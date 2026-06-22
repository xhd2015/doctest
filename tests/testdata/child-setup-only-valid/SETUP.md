# Scenario

**Feature**: the root defines Request, Response, and Run

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Child Setup Only Valid Fixture

## Preconditions
- The root defines Request, Response, and Run.
- An intermediate child defines only Setup (no Run), which is valid.
- The leaf verifies the combined behavior runs correctly.