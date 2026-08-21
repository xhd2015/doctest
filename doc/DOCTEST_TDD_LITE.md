---
name: doctest-tdd-lite
description: single-agent doctest TDD (design + implement inline, no sub-agent delegation)
---

--begin of skill doctest-tdd-lite--

# Gate

TDD is mandatory **unless**: (1) `no tdd:` prefix — strip it and follow
`doctest-dev-test`; (2) doc-only (`.md`, README, etc.); (3) one-liner fix — warn
TDD is slow, ask TDD vs `no tdd:`, do not start without confirmation.

# One Rule

You are a **single agent** running doctest TDD. You write the doctest tree
**and** the implementation. Every code change flows through:

```
Classic TDD:       design tests → RED → seal → implement → GREEN
Coverage backfill: design tests → (GREEN OK / mixed OK) → seal → implement only if RED remains → GREEN
```

You may use a sub-agent only for explore/analysis. All file edits (tests and
source) are yours. Plan phases in context → one full cycle **per phase**
(see **Plan phases**), not one mega-cycle.

# Modes

Infer in TDD step 1. Do not force classic RED when implementation is present.

| Mode | When |
|------|------|
| **Classic TDD** | Greenfield / stubbed / user wants failing tests first |
| **Coverage backfill** | User asks to **backfill**, or code already implements intended behavior and only doctests are missing |

Ambiguous → ask once; default backfill when evidence is strong. Brainstorm in
**both** modes (step 1 mandatory for backfill too).

**Shortcuts:** bug already RED → jump to seal. Backfill: RED not required for
correct behavior (GREEN expected); skip implement if all sealed tests GREEN.

**Backfill (design)** — when writing tests in backfill mode: assert existing
correct behavior; GREEN expected; mixed GREEN/RED OK when some behaviors still
missing; no must-fail asserts only for classic TDD theater.

# Test layers (L1 / L2 / L3)

Sole definition here (guidance, not CI). Deep dive:
`doctest skill design-principle --show`.

| Layer | ~Cases | Execution | Use for |
|-------|--------|-----------|---------|
| **L1** | 10–20% | `*_test.go` tables | Pure / flat edges |
| **L2** ★ | **70–85%** | Same process: library or `cli.RunWithWriter` | Multi-factor APIs + short CLI |
| **L3** | 5–10% | Separate process / nested suite; **`label: e2e`** | Full integration only |

Default **L2**. Short path → never L3. L3 needs process-boundary reason +
`label: e2e`. “Prefer doctests” = prefer L2, not e2e. Requirements: state a
**layer map** only — do not restate this table.

# Test-first code

Design **easy to test**, not only production-correct: injectable L2 APIs (opts,
writers, pure cores); no GREEN via Setenv/Chdir/stdio or forced L3 for short
paths (see **Parallel-safe suite**, **Test layers**); e2e-only / untestable
glue GREEN is incomplete. Design & implement: apply; do not restate.

# Parallel-safe suite

Leaves use `t.Parallel()` in one process. **Forbidden** in harness and L2
product: (1) unprotected shared globals / reassigning `os.Stdout|Stderr|Stdin`;
(2) process-global env/cwd — `os.Setenv`/`Unsetenv`, `os.Chdir`, `t.Setenv`,
`t.Chdir`, `syscall.Setenv`. **Prefer** inject opts / `req` / child
`cmd.Env`·`Dir`. Detail: `doctest skill lint --show`, `doctest skill review
--show` (Common gotchas). Design & implement: apply; do not restate.

# Plan phases (outer loop)

**Plan phase** = dependency-ordered unit (`P1`… / split-phases / `PHASES.md`).
**TDD step** = inner workflow step 1–6 — do not confuse the two.

**Trigger:** plan phases in context, or user asks phase-by-phase. **None:**
single cycle (files omit `PHASE` — see **Requirement naming**).

```text
for each Pn in order (or user subset):
  scope to Pn → TDD steps 1→6 → on GREEN auto-continue → stop when done
```

| Hard rule | |
|-----------|---|
| One phase = one full TDD cycle | No design-all-then-implement-all |
| Scope to phase exit criteria | Later work only if plan allows stubs/seams |
| Mode per phase | Classic vs backfill may differ |
| Paths / trees | **Requirement naming**; `./tests/<feature>/` or `<pkg>/tests/<feature>/` (no phase subdirs) |
| Seal once per cycle | No unjustified rewrite of prior sealed asserts |
| Auto-continue | Then short multi-phase summary |

# Workflow (6 TDD steps)

## TDD step 1 — Requirements

Brainstorm both modes; write requirement; get approval. Mode per **Modes**.
Plan phase active → scope `Pn` (**Plan phases**). Tell user: (1) models/layout
if any (2) scenarios + expected output (3) layer map (**Test layers**)
(4) classic vs backfill + why (5) `Pn` or single-cycle. Bugs: explore then
write doctests; classic → RED; backfill (fix applied) → GREEN OK.

## TDD step 2 — Design & Write Tests

Design and materialize the tree under `<pkg>/tests/<feature>/` (or
`./tests/<feature>/`). Follow **DOCTEST_SPEC** and **DOCTEST_DESIGN_SPEC**
below (MECE, Setup/Assert, DSN, session, output asserts). Apply **Test
layers**, **Test-first code**, **Parallel-safe suite**. Existing tree → find
gaps and **backfill** or update. Backfill rules: see **Modes**. Multi-phase →
scope leaves to `Pn` exit criteria.

## TDD step 3 — Vet then run

```sh
doctest vet ./tests/<feature>
doctest test ./tests/<feature>
```

`vet` must pass (`__DOCTEST_VERSION__` — SPEC below). Interpret test per
**Modes**: classic → RED (any GREEN → re-examine); backfill → GREEN expected
(re-examine vacuous/wrong only); mixed → seal as-is.

## TDD step 4 — Seal (once)

```
git add ./tests/<feature>
```

Tests only, once per cycle. No git repo → ask before unsealed. Mixed seal as-is.

## TDD step 5 — Implement

Backfill + all sealed GREEN → **skip** to step 6; report production code
unchanged and doctests backfilled only.

Otherwise implement until sealed tests pass:

- **Never modify sealed test files** — if an assert seems wrong, ask the user
- Do not weaken already-GREEN sealed asserts
- Implementation in source (not `_test.go` for harness logic)
- Match `Request` / `Response` / `Run` from root `DOCTEST.md`
- Paths/session only via **`d`**: `d.DOCTEST_ROOT` / `d.DOCTEST_CASE` /
  `d.DOCTEST_SESSION_ID` (not free vars or getenv)
- Apply **Test-first code** and **Parallel-safe suite** (pointers only)

```sh
doctest test ./tests/<feature>
doctest test ./...
```

## TDD step 6 — Verify

```sh
git diff ./tests/<feature>            # must be clean
doctest vet ./tests/<feature>
doctest test ./tests/<feature>/...    # must be GREEN
doctest test --label "ui-automation" ./tests/<feature>/... # if ASSERT has label
doctest test --label 'slow && ui-automation' ./tests/<feature>/... # &&, ||, ()
doctest test ./...                    # no regressions
```

RED → fix implementation (step 5). Test-file edits only with explicit user
justification. Report count + accepted mods. More plan phases → **Plan
phases**; else summary.

# Requirement naming

Under `/tmp/` (not repo root):

- single-cycle: `/tmp/REQUIREMENT-<slug>.md`
- plan phase N: `/tmp/REQUIREMENT-PHASE-{N}-<slug>.md`

`PHASE-{N}` only when multi-phase — never invent for single-cycle.

**Must include** (not restated in steps): context/summary; mode (+ **Backfill
(design)** from **Modes** if backfill); layer map (**Test layers**); apply
**Test-first code** + **Parallel-safe suite** (pointers); scenarios / planned
tree; multi-phase: `Pn` goal/exit/out-of-scope. After seal: **"tests are
sealed — do not modify"** and the verify command.

# Followup

Restart at TDD step 1 (relevant plan phase if multi-phase). Ask the user
directly for clarification in the normal response (no parent orchestrator).

# SPECS

<DOCTEST_SPEC>
__DOCTEST_SPEC__
</DOCTEST_SPEC>

<DOCTEST_DESIGN_SPEC>
__DOCTEST_DESIGN_SPEC__
</DOCTEST_DESIGN_SPEC>

--end of skill doctest-tdd-lite--
