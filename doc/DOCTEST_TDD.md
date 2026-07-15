---
name: doctest-tdd
description: adversarial multi-agent TDD with doctests (orchestrator + tests designer + implementer)
---

--begin of skill doctest-tdd--

# Gate

TDD mode is mandatory **unless** one of these applies:

1. **`no tdd:` prefix** — Strip the prefix and handle directly.
2. **Doc-only change** — Only documentation (`.md`, `README`, etc.); no source code.
3. **One-liner fix** — Warn that TDD is slow; ask whether to proceed with TDD or use
   `no tdd:`. Do not start TDD without confirmation.

Otherwise TDD mode is mandatory.

# One Rule

You are the **orchestrator**. Every code change MUST flow through doctest TDD:

```
Classic TDD:       designer → RED → seal → implementer → GREEN
Coverage backfill: designer → (GREEN OK / mixed OK) → seal → implementer only if RED remains → GREEN
```

When the context already has a **plan split into phases** (see **Plan phases**
below), apply this rule **once per plan phase** — full TDD cycle each time —
not one mega-cycle for the whole plan.

**You NEVER touch source files.** No Edit/Write on source or config. All code
changes go through:

- **Designer** — writes doctest trees
- **Implementer** — writes implementation code

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

- Already-reproduced bug with RED doctests → skip design/RED; start at seal.
- Backfill: RED not required for leaves documenting correct behavior; GREEN
  expected. Skip implementer when all sealed tests are already GREEN.

# Plan phases (outer loop)

**Plan phase** = a dependency-ordered work unit from a split plan (`P1`, `P2`,
…; split-phases output; `PHASES.md`; or an equivalent phase list in context).

**TDD step** = one step of the inner workflow below (TDD steps 1–8). Do not
confuse the two.

**Trigger** — any of:

- Context has a split-phases (or equivalent) plan with plan phases and exit criteria
- `PHASES.md` (or similar) is the agreed plan
- User asks to implement phase-by-phase / by plan phase

**If no plan phases in context:** run a single inner TDD cycle as today
(requirement files still use `_PHASE_1`).

**If plan phases are present:**

```text
for each plan phase Pn in dependency order
    (or only the subset the user named):
  1. Scope requirements to Pn only (goal, work, exit criteria, out of scope)
  2. Run full TDD steps 1→8 for Pn
  3. On GREEN + verify: auto-continue to the next plan phase
  4. Stop only when all plan phases in scope are done
```

**Hard rules:**

1. **One plan phase = one full TDD cycle** — do not design all phases’ tests
   then implement everything in one pass.
2. **Scope to that phase’s exit criteria** — do not pull later plan-phase work
   forward (stubs/seams OK only if the phase plan allows).
3. **Mode per plan phase** — classic vs backfill may differ by phase.
4. **Requirement files** use `REQUIREMENT_DESIGN_PHASE_n` /
   `REQUIREMENT_IMPLEMENT_PHASE_n` (see naming below).
5. **Doctest tree paths stay normal** (`./tests/<feature>/`) — no required
   phase subdirs; later plan phases may add leaves under the same tree.
6. **Seal once per TDD cycle** (i.e. per plan phase when multi-phase) — seal
   that cycle’s new/changed tests; do not rewrite prior sealed asserts without
   justification.
7. **Auto-continue** until every in-scope plan phase is done; then report a
   short summary across phases.

# Delegating to roles

Spawn/resume role subagents via the runner's task tool. Pass a short requirement
(or file path) — not the role prompt.

Each sub-agent **must** run as its first step:

- Designer: `doctest skill designer --show`
- Implementer: `doctest skill implementer --show`

Prefix descriptions with `[designer]` / `[implementer]` when supported. Resume
the **same** session per role for follow-ups. Wait patiently (use ≥1h timeout
if the runner requires one).

# Workflow (8 TDD steps)

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

CLI: user-facing stdout ends with `\n` after the last content line; when using
doctest assert templates, newline before the raw string's closing backtick.

Bugs: explore/narrow scope first, then designer. Classic → failing doctests;
backfill (fix applied) → doctests of fixed behavior (GREEN OK).

Clarify until intent is clear.

## TDD step 2 — Delegate Test Design

Write `REQUIREMENT_DESIGN_PHASE_<n>.md`. Spawn designer with that path (+ optional
summary). Designer runs `doctest skill designer --show` first.

**Backfill — MUST tell the designer** (spawn message and/or requirement):

- Mode: coverage backfill (implementation present)
- Intent: **backfill** missing doctests for existing correct behavior
- RED not required; GREEN expected for covered paths
- Mixed GREEN/RED OK when some behaviors still missing
- Do not invent must-fail assertions only for classic TDD theater

Wait until the tree is under `./tests/<feature>/`.

## TDD step 3 — Designer Questions (optional)

Answer questions by **resuming the same designer session**. Escalate
domain questions to the user first. Repeat until designer completes.

## TDD step 4 — Vet then run

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

## TDD step 5 — Seal (once)

```
git add ./tests/<feature>
```

Seal tests only, once per TDD cycle. Outside a git repo, ask before proceeding
unsealed. **Mixed suites seal as-is** (GREEN + RED together).

## TDD step 6 — Implement

If backfill and all sealed tests are GREEN → **skip** (no implementer); go to
TDD step 8.

Else write `REQUIREMENT_IMPLEMENT_PHASE_<n>.md`: context, summary, tree structure,
**"tests are sealed — do not modify"**, verify command, which leaves were
already GREEN vs still RED. Include active plan phase when multi-phase.

Spawn implementer with that path. Implementer runs `doctest skill implementer
show` first. Wait until all tests pass. Do not weaken already-GREEN sealed
asserts.

## TDD step 7 — Implementer Questions (optional)

Resume the **same implementer session** with answers or TDD step 8 failures until
all pass.

## TDD step 8 — Verify

```sh
git diff ./tests/<feature>            # must be clean
doctest vet ./tests/<feature>
doctest test ./tests/<feature>/...    # must be GREEN

doctest test --label "ui-automation" ./tests/<feature>/... # if ASSERT has label header
doctest test --label 'slow && ui-automation' ./tests/<feature>/... # expr: &&, ||, ()

doctest test ./...                    # no regressions
```

If RED → resume implementer. Accept test-file changes only with explicit
justification (wrong expected per spec). Report test count and any accepted
modifications.

When plan phases remain in scope: auto-continue to the next plan phase (new
requirement files, TDD steps 1→8 again). When all in-scope plan phases are
done: report a short multi-phase summary.

# Requirement File Naming

- Design: `REQUIREMENT_DESIGN_PHASE_<n>.md` — must state mode; backfill includes
  TDD step 2 handoff bullets; multi-phase runs state plan phase `Pn` scope
- Implement: `REQUIREMENT_IMPLEMENT_PHASE_<n>.md`

Use `n` = plan phase number when multi-phase; use `1` for a single-cycle run
with no plan split.

# Followup Requests

Every followup restarts at TDD step 1 (for the relevant plan phase if still
multi-phase). Keep designer and implementer sessions separate; reuse each
role's session ID within that role.

__DOCTEST_SPEC__

--end of skill doctest-tdd--
