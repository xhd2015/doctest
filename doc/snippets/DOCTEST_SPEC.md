# Doctest specification

A doc-style test is a test expressed as markdown in a decision tree.
Prose is primary, code supplementary.

## Layout

```
<pkg>/tests/<feature>/
├── DOCTEST.md          # Overview, diagram, test index + "## How to Run"
├── SETUP.md            # Root: shared preconditions, Request/Response types, the Run(t, req) (resp,error) function that actually runs the logic; can use DOCTEST_ROOT to refer this dir
│                       #   Request, Response, and Run cannot be redefined by any descendant
├── <decision-slug>/           # Grouping — no ASSERT.md, must have SETUP.md
│   └── <leaf-slug>/           # Runnable — has ASSERT.md, must have SETUP.md
│       └── testdata/   # Fixtures (skipped)
```

- Dir with `ASSERT.md` = runnable leaf; without = grouping node (must have SETUP.md)
- `SETUP.md` accumulates root→leaf: `## Steps` concatenated, `## Preconditions`/`## Context` merged.
  Sections: `## Preconditions`, `## Steps`, `## Context`
- `ASSERT.md` is **case-private** (never inherited).
  Sections: `## Expected`, `## Side Effects`, `## Errors`, `## Exit Code`

## Vet

Run `doctest vet <dir>` to validate tree structure before `doctest test` or
`doctest build`. Vet checks:

- Root `DOCTEST.md` exists and includes a **DSN (Domain Specific Notion)** section
  (heading text must contain `DSN (Domain Specific Notion)`)
- Every `SETUP.md` starts with `# Scenario` as its first section (leading
  whitespace allowed)
- `ASSERT.md` directories also have `SETUP.md`
- Anti-patterns in `SETUP.md` / `ASSERT.md` code blocks (embedded Go programs,
  `go test` shell-outs)

See `doctest skill doc-spec show` and `doctest skill code-spec show` for the
full prose and code rules; see the design spec (`doctest skill designer show`)
for DSN and Scenario semantics.
