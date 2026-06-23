---
name: doc-style-test-based-tdd-lite
description: single-agent doctest TDD (design + implement inline, no sub-agent delegation)
---

--begin of skill doc-style-test-based-tdd-lite--

# Gate

TDD mode is mandatory **unless** one of these applies:

1. **`no tdd:` prefix** — The user's message starts with `no tdd:`. Strip the
   prefix and handle directly.
2. **Doc-only change** — The change touches only documentation files (`.md`,
   `README`, etc.) and modifies no source code. No TDD flow needed.
3. **One-liner fix** — Before starting the expensive TDD loop, ask the user to
   confirm: warn that TDD is slow and ask whether they want to proceed
   with TDD or use `no tdd:` instead. Do not start TDD without confirmation.

Otherwise TDD mode is mandatory. Everything below is non-negotiable.

# TDD_LITE Gate — The One Rule

You are a **single agent** running doctest TDD. Every code change — features,
bug fixes, trivial one-liners, variable renames, logging changes, config tweaks
— MUST flow through:

```
design tests → RED → seal → implement → GREEN
```

You write the doctest tree **and** the implementation yourself.

If you've already reproduced a bug with written doctests, and confirm they're RED, then you can jump start from seal step. 

You may use sub-agent to explore/analysis to narrow bug scope before writing
tests. All file edits (tests and source) are yours.

# Workflow (6 Phases)

Every feature request, bug fix, or followup follows this loop.

## Phase 1 — Requirements

Brainstorm with user. Produce a requirement file. Get explicit user approval
before continuing.

Explicitly tell user:

1. What the underlying data models and storage layout (if any) are;
2. What scenarios you will test, and expected output;
3. How you gonna test that, prefer rerunable tests (doc-style tests or unit tests);

For bugs: reproduce is critical. Use read-only explore/analysis if needed to
narrow scope, then write failing doctests that capture the bug.

If there is anything needs clarification, list them and ask for user
confirmation until no gap understanding user's intent.

## Phase 2 — Design & Write Tests

Design and materialize the doctest tree yourself. Follow the doc-style test
specification(see below DOCTEST_SPEC and DOCTEST_DESIGN_SPEC)

### Understand the requirement

Identify inputs, outputs, side effects, every flag/option/mode, happy paths,
error paths, and edge cases.

### Design the decision tree

- Split on the **most significant** parameter at the root (e.g. operation mode)
- Each grouping level narrows on exactly one parameter with **mutually
  exclusive** branches
- Recurse until concrete runnable leaves
- Cover every valid path and every error path

If a relevant doctest tree already exists, inspect it, find coverage gaps, and
add or update tests.

### Write the tree

Materialize under `<pkg>/tests/<feature>/`:

- Root `DOCTEST.md`: DSN, `## Version` (`__DOCTEST_VERSION__`), ASCII decision
  tree diagram, test-leaf index, `## How to Run`, and the Go block defining
  `Request`, `Response`, and `Run`
- Grouping nodes: `SETUP.md` only (no `ASSERT.md`)
- Leaves: `SETUP.md` + `ASSERT.md`

Generated tests expose `DOCTEST_ROOT` and `DOCTEST_SESSION_ID`. Use
`DOCTEST_SESSION_ID` for session-scoped shared directories or locks.

Coverage checklist — every leaf should cover:

- Happy paths for valid input combinations
- Error paths for invalid input
- Edge cases (empty, zero, boundary, extremes)
- Parameter interactions
- Prefer more leaves over fewer

If the feature is already correctly implemented and tests pass before any code
change, report that to the user instead of unnecessary implementation.

## Phase 3 — Vet then RED

```sh
doctest vet ./tests/<feature>
doctest test ./tests/<feature>
```

`doctest vet` must pass — the tree is well-formed (including `## Version` and
DSN in root `DOCTEST.md`, Request/Response/Run in the `DOCTEST.md` Go block,
and `# Scenario` as the first section in every `SETUP.md`). Current spec
version: `__DOCTEST_VERSION__`.

`doctest test` should fail (RED) since no implementation exists yet. If any
test passes, re-examine the test design.

## Phase 4 — Seal (once)

```
git add ./tests/<feature>
```

Never seal more than once. Only tests get sealed, never code. If outside a git
repo, ask the user before proceeding unsealed.

## Phase 5 — Implement

Write implementation code to make all sealed tests pass. Follow these rules:

- **Never modify sealed test files** — if an assertion seems wrong, ask the
  user for clarification rather than editing tests
- Place implementation in appropriate source files (not `_test.go` for doctest
  harness logic)
- Use types and signatures expected by the root `DOCTEST.md` Go block
- Use `DOCTEST_SESSION_ID` for per-run session-scoped cache paths or
  coordination

Run tests until GREEN:

```sh
doctest test ./tests/<feature>
```

If any test fails, fix implementation and re-run. Also run full regression:

```sh
doctest test ./...
```

## Phase 6 — Verify

```sh
git diff ./tests/<feature>            # must be clean
doctest vet ./tests/<feature>         # structure still valid
doctest test ./tests/<feature>/...    # must be GREEN
doctest test ./...                    # no regressions
```

If RED, fix implementation and repeat Phase 5. If test files were modified,
accept only with explicit user justification (the test expected wrong behavior
per spec).

Report: test count, modifications accepted (with rationale).

# Requirement File Naming

Use a single file per feature:

- `REQUIREMENT-<slug>.md`

Include:

- Summarized context and feature summary
- Data models and storage layout (if any)
- Test scenarios and expected outputs
- Planned test tree structure
- After sealing: **"tests are sealed — do not modify"** and the verify command

Use `--requirement` when invoking tools that accept it, or keep the file in the
workspace for your own reference during the loop.

# Followup Requests

Every followup (new feature, fix, amendment) restarts this workflow from
Phase 1.

Ask the user directly when you need clarification. Do not use
`yield-pending-questions` or sub-agent session IDs — there is no parent
orchestrator.

# SPECS

<DOCTEST_SPEC>
__DOCTEST_SPEC__
</DOCTEST_SPEC>

<DOCTEST_DESIGN_SPEC>
__DOCTEST_DESIGN_SPEC__
</DOCTEST_DESIGN_SPEC>

--end of skill doc-style-test-based-tdd-lite--