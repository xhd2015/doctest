---
name: doctest-tdd-lite
description: single-agent doctest TDD (design + implement inline, no sub-agent delegation)
---

--begin of skill doctest-tdd-lite--

# Gate

TDD mode is mandatory **unless** one of these applies:

1. **`no tdd:` prefix** — Strip the prefix and handle directly.
2. **Doc-only change** — Only documentation (`.md`, `README`, etc.); no source code.
3. **One-liner fix** — Warn that TDD is slow; ask whether to proceed with TDD or use
   `no tdd:`. Do not start TDD without confirmation.

Otherwise TDD mode is mandatory.

# One Rule

You are a **single agent** running doctest TDD. You write the doctest tree
**and** the implementation. Every code change MUST flow through:

```
Classic TDD:       design tests → RED → seal → implement → GREEN
Coverage backfill: design tests → (GREEN OK / mixed OK) → seal → implement only if RED remains → GREEN
```

You may use a sub-agent only for explore/analysis. All file edits (tests and
source) are yours.

# Modes

Infer mode from user wording and codebase during Phase 1. Do not force classic
RED when the implementation is already present.

**Classic TDD** when: greenfield / no real implementation; user wants failing
tests first; paths under change are stubbed or unimplemented.

**Coverage backfill** when: user asks to **backfill** coverage or says the
fix/feature is already applied and only doctests are missing; source already
implements the intended behavior; task is missing coverage for working code.

If ambiguous, ask once in Phase 1 — default toward backfill when evidence is strong.

**Still brainstorm in both modes** (Phase 1 is mandatory for backfill too).

### Shortcuts

- Already-reproduced bug with RED doctests → jump to seal.
- Backfill: RED not required for leaves documenting correct behavior; GREEN
  expected. Skip implement when all sealed tests are already GREEN.

# Workflow (6 Phases)

## Phase 1 — Requirements

Brainstorm (both modes). Produce a requirement file; get explicit approval.

Auto-detect mode (see **Modes**). State it in the requirement file.

Tell the user:

1. Data models and storage layout (if any)
2. Scenarios and expected output
3. How you will test (prefer doctests)
4. **Classic TDD** or **coverage backfill**, and why

Bugs: explore if needed, then write doctests. Classic → failing doctests;
backfill (fix applied) → **backfill** doctests of fixed behavior (GREEN OK).

Clarify until intent is clear.

## Phase 2 — Design & Write Tests

Design and materialize the doctest tree yourself under
`<pkg>/tests/<feature>/`. Follow **DOCTEST_SPEC** and **DOCTEST_DESIGN_SPEC**
below (MECE tree, Setup/Assert, DSN, session setup, output asserts).

- If a tree already exists, find coverage gaps and **backfill** or update tests.
- **Backfill mode:** assert existing correct behavior; GREEN expected; mixed
  GREEN/RED OK when some behaviors still missing; do not invent must-fail
  asserts only for classic TDD theater.

## Phase 3 — Vet then run

```sh
doctest vet ./tests/<feature>
doctest test ./tests/<feature>
```

`doctest vet` must pass (tree well-formed; version `__DOCTEST_VERSION__` — see
SPEC below). Interpret `doctest test` per **Modes**:

- **Classic:** must be RED; any GREEN → re-examine design.
- **Backfill:** GREEN expected for covered behavior; re-examine only vacuous
  or wrong asserts.
- **Mixed:** valid; seal the whole tree as-is.

## Phase 4 — Seal (once)

```
git add ./tests/<feature>
```

Seal tests only, once. Outside a git repo, ask before proceeding unsealed.
**Mixed suites seal as-is** (GREEN + RED together).

## Phase 5 — Implement

If backfill and all sealed tests are GREEN → **skip**; go to Phase 6. Report
that production code was unchanged and doctests were backfilled only.

Otherwise implement until sealed tests pass:

- **Never modify sealed test files** — if an assert seems wrong, ask the user
- Do not weaken already-GREEN sealed asserts
- Put implementation in source (not `_test.go` for harness logic)
- Match `Request` / `Response` / `Run` from root `DOCTEST.md`
- Use injected `DOCTEST_SESSION_ID` (not `os.Getenv`)

```sh
doctest test ./tests/<feature>
doctest test ./...
```

## Phase 6 — Verify

```sh
git diff ./tests/<feature>            # must be clean
doctest vet ./tests/<feature>
doctest test ./tests/<feature>/...    # must be GREEN

doctest test --label "ui-automation" ./tests/<feature>/... # if ASSERT has label header
doctest test --label 'slow && ui-automation' ./tests/<feature>/... # expr: &&, ||, ()

doctest test ./...                    # no regressions
```

If RED → fix implementation (Phase 5). Accept test-file changes only with
explicit user justification. Report test count and any accepted modifications.

# Requirement File Naming

Single file: `REQUIREMENT-<slug>.md`

Include: context/summary, mode (classic or backfill + why), data models,
scenarios, planned tree. After seal: **"tests are sealed — do not modify"**
and the verify command.

# Followup Requests

Every followup restarts at Phase 1. Ask the user directly for clarification
(no parent orchestrator / yield-pending-questions).

# SPECS

<DOCTEST_SPEC>
__DOCTEST_SPEC__
</DOCTEST_SPEC>

<DOCTEST_DESIGN_SPEC>
__DOCTEST_DESIGN_SPEC__
</DOCTEST_DESIGN_SPEC>

--end of skill doctest-tdd-lite--
