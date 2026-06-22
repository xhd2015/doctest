# Doctest specification

A doc-style test is a test expressed as markdown in a decision tree.
Prose is primary, code supplementary.

## Layout

```
<pkg>/tests/<feature>/
├── DOCTEST.md          # Overview, ## Version (__DOCTEST_VERSION__), diagram, test index; final ```go``` block defines Request/Response/Run
├── SETUP.md            # Optional root: func Setup + helper funcs shared by all tests (no Request/Response/Run)
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

- Root `DOCTEST.md` exists and includes **DSN (Domain Specific Notion)** and
  **## Version** sections (version presence only; value is `__DOCTEST_VERSION__`)
- Root `DOCTEST.md` Go block defines `type Request`, `type Response`, and `func Run`
- Root `SETUP.md` (when present) must not define Request, Response, or Run
- Every `SETUP.md` starts with `# Scenario` as its first section (leading
  whitespace allowed)
- `ASSERT.md` directories also have `SETUP.md`
- Anti-patterns in `SETUP.md` / `ASSERT.md` code blocks (embedded Go programs,
  `go test` shell-outs)

See `doctest skill doc-spec show` and `doctest skill code-spec show` for the
full prose and code rules; see the design spec (`doctest skill designer show`)
for DSN and Scenario semantics.
