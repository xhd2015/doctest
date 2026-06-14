# Doctest specification

A doc-style test is a test expressed as markdown in a decision tree.
Prose is primary, code supplementary.

## Layout

```
<pkg>/tests/<feature>/
├── DOCTEST.md          # Overview, diagram, test index + "## How to Run"
├── SETUP.md            # Root: shared preconditions, Request/Response types, the Run(t, req) (resp,error) function that actually runs the logic;can use DOCTEST_ROOT to refer this dir
├── decision/           # Grouping — no ASSERT.md, must have SETUP.md
│   └── leaf/           # Runnable — has ASSERT.md, must have SETUP.md
│       └── testdata/   # Fixtures (skipped)
```

- Dir with `ASSERT.md` = runnable leaf; without = grouping node (must have SETUP.md)
- `SETUP.md` accumulates root→leaf: `## Steps` concatenated, `## Preconditions`/`## Context` merged.
  Sections: `## Preconditions`, `## Steps`, `## Context`
- `ASSERT.md` is **case-private** (never inherited).
  Sections: `## Expected`, `## Side Effects`, `## Errors`, `## Exit Code`
