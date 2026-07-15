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

When the context already has a **plan split into phases** (see **Plan phases**
below), apply this rule **once per plan phase** — full TDD cycle each time —
not one mega-cycle for the whole plan.

You may use a sub-agent only for explore/analysis. All file edits (tests and
source) are yours.

# Modes

Infer mode from user wording and codebase during TDD step 1. Do not force classic
RED when the implementation is already present.

**Classic TDD** when: greenfield / no real implementation; user wants failing
tests first; paths under change are stubbed or unimplemented.

**Coverage backfill** when: user asks to **backfill** coverage or says the
fix/feature is already applied and only doctests are missing; source already
implements the intended behavior; task is missing coverage for working code.

If ambiguous, ask once in TDD step 1 — default toward backfill when evidence is strong.

**Still brainstorm in both modes** (TDD step 1 is mandatory for backfill too).

### Shortcuts

- Already-reproduced bug with RED doctests → jump to seal.
- Backfill: RED not required for leaves documenting correct behavior; GREEN
  expected. Skip implement when all sealed tests are already GREEN.

# Plan phases (outer loop)

**Plan phase** = a dependency-ordered work unit from a split plan (`P1`, `P2`,
…; split-phases output; `PHASES.md`; or an equivalent phase list in context).

**TDD step** = one step of the inner workflow below (TDD steps 1–6). Do not
confuse the two.

**Trigger** — any of:

- Context has a split-phases (or equivalent) plan with plan phases and exit criteria
- `PHASES.md` (or similar) is the agreed plan
- User asks to implement phase-by-phase / by plan phase

**If no plan phases in context:** run a single inner TDD cycle as today
(requirement file still uses `_PHASE_1`).

**If plan phases are present:**

```text
for each plan phase Pn in dependency order
    (or only the subset the user named):
  1. Scope requirements to Pn only (goal, work, exit criteria, out of scope)
  2. Run full TDD steps 1→6 for Pn
  3. On GREEN + verify: auto-continue to the next plan phase
  4. Stop only when all plan phases in scope are done
```

**Hard rules:**

1. **One plan phase = one full TDD cycle** — do not design all phases’ tests
   then implement everything in one pass.
2. **Scope to that phase’s exit criteria** — do not pull later plan-phase work
   forward (stubs/seams OK only if the phase plan allows).
3. **Mode per plan phase** — classic vs backfill may differ by phase.
4. **Requirement files** use `REQUIREMENT_PHASE_n` (see naming below).
5. **Doctest tree paths stay normal** (`./tests/<feature>/` or
   `<pkg>/tests/<feature>/`) — no required phase subdirs; later plan phases
   may add leaves under the same tree.
6. **Seal once per TDD cycle** (i.e. per plan phase when multi-phase) — seal
   that cycle’s new/changed tests; do not rewrite prior sealed asserts without
   justification.
7. **Auto-continue** until every in-scope plan phase is done; then report a
   short summary across phases.

# Workflow (6 TDD steps)

## TDD step 1 — Requirements

Brainstorm (both modes). Produce a requirement file; get explicit approval.

Auto-detect mode (see **Modes**). State it in the requirement file.

When a plan phase is active: state which plan phase (`Pn`), its goal, exit
criteria, and out of scope; scope scenarios to that phase only.

Tell the user:

1. Data models and storage layout (if any)
2. Scenarios and expected output
3. How you will test (prefer doctests)
4. **Classic TDD** or **coverage backfill**, and why
5. Active plan phase (`Pn`) when multi-phase, or that this is a single-cycle run

Bugs: explore if needed, then write doctests. Classic → failing doctests;
backfill (fix applied) → **backfill** doctests of fixed behavior (GREEN OK).

Clarify until intent is clear.

## TDD step 2 — Design & Write Tests

Design and materialize the doctest tree yourself under
`<pkg>/tests/<feature>/`. Follow **DOCTEST_SPEC** and **DOCTEST_DESIGN_SPEC**
below (MECE tree, Setup/Assert, DSN, session setup, output asserts).

- If a tree already exists, find coverage gaps and **backfill** or update tests.
- **Backfill mode:** assert existing correct behavior; GREEN expected; mixed
  GREEN/RED OK when some behaviors still missing; do not invent must-fail
  asserts only for classic TDD theater.
- Scope new leaves to the active plan phase’s exit criteria when multi-phase.

## TDD step 3 — Vet then run

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

## TDD step 4 — Seal (once)

```
git add ./tests/<feature>
```

Seal tests only, once per TDD cycle. Outside a git repo, ask before proceeding
unsealed. **Mixed suites seal as-is** (GREEN + RED together).

## TDD step 5 — Implement

If backfill and all sealed tests are GREEN → **skip**; go to TDD step 6. Report
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

## TDD step 6 — Verify

```sh
git diff ./tests/<feature>            # must be clean
doctest vet ./tests/<feature>
doctest test ./tests/<feature>/...    # must be GREEN

doctest test --label "ui-automation" ./tests/<feature>/... # if ASSERT has label header
doctest test --label 'slow && ui-automation' ./tests/<feature>/... # expr: &&, ||, ()

doctest test ./...                    # no regressions
```

If RED → fix implementation (TDD step 5). Accept test-file changes only with
explicit user justification. Report test count and any accepted modifications.

When plan phases remain in scope: auto-continue to the next plan phase (new
requirement file, TDD steps 1→6 again). When all in-scope plan phases are
done: report a short multi-phase summary.

# Requirement File Naming

Single file: `REQUIREMENT_PHASE_<n>.md`

Use `n` = plan phase number when multi-phase; use `1` for a single-cycle run
with no plan split.

Include: context/summary, mode (classic or backfill + why), data models,
scenarios, planned tree; when multi-phase, plan phase `Pn` goal/exit
criteria/out of scope. After seal: **"tests are sealed — do not modify"**
and the verify command.

# Followup Requests

Every followup restarts at TDD step 1 (for the relevant plan phase if still
multi-phase). Ask the user directly for clarification
(no parent orchestrator / yield-pending-questions).

# SPECS

<DOCTEST_SPEC>
__DOCTEST_SPEC__
</DOCTEST_SPEC>

<DOCTEST_DESIGN_SPEC>
__DOCTEST_DESIGN_SPEC__
</DOCTEST_DESIGN_SPEC>

--end of skill doctest-tdd-lite--
