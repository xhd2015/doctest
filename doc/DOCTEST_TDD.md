---
name: doctest-tdd
description: adversarial multi-agent TDD with doctests (orchestrator + tests designer + implementer)
---

--begin of skill doctest-tdd--

# Gate

TDD is mandatory **unless**: (1) `no tdd:` prefix — strip it and follow
`doctest-dev-test`; (2) doc-only (`.md`, README, etc.); (3) one-liner fix — warn
TDD is slow, ask TDD vs `no tdd:`, do not start without confirmation.

# One Rule

You are the **orchestrator**. Every code change flows through:

```
Classic TDD:       designer → RED → seal → implementer → GREEN
Coverage backfill: designer → (GREEN OK / mixed OK) → seal → implementer only if RED remains → GREEN
```

**You NEVER touch source/config.** Designer writes doctest trees; implementer
writes implementation. Plan phases in context → one full cycle **per phase**
(see **Plan phases**), not one mega-cycle.

# Modes

Infer in TDD step 1. Do not force classic RED when implementation is present.

| Mode | When |
|------|------|
| **Classic TDD** | Greenfield / stubbed / user wants failing tests first |
| **Coverage backfill** | User asks to **backfill**, or code already implements intended behavior and only doctests are missing |

Ambiguous → ask once; default backfill when evidence is strong. Brainstorm in
**both** modes (step 1 mandatory for backfill too).

**Shortcuts:** bug already RED → start at seal. Backfill: RED not required for
correct behavior (GREEN expected); skip implementer if all sealed tests GREEN.

**Backfill handoff (designer)** — REQUIREMENT-DESIGN / spawn must include:
mode = coverage backfill; intent = backfill missing doctests for correct
behavior; RED not required / GREEN expected; mixed GREEN/RED OK; no must-fail
asserts only for classic TDD theater.

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
glue GREEN is incomplete. Handoffs: apply; do not restate.

# Parallel-safe suite

Leaves use `t.Parallel()` in one process. **Forbidden** in harness and L2
product: (1) unprotected shared globals / reassigning `os.Stdout|Stderr|Stdin`;
(2) process-global env/cwd — `os.Setenv`/`Unsetenv`, `os.Chdir`, `t.Setenv`,
`t.Chdir`, `syscall.Setenv`. **Prefer** inject opts / `req` / child
`cmd.Env`·`Dir`. Detail: `doctest skill lint --show`, `doctest skill review
--show` (Common gotchas). Handoffs: apply; do not restate.

# Plan phases (outer loop)

**Plan phase** = dependency-ordered unit (`P1`… / split-phases / `PHASES.md`).
**TDD step** = inner workflow step 1–8 — do not confuse the two.

**Trigger:** plan phases in context, or user asks phase-by-phase. **None:**
single cycle (files omit `PHASE` — see **Requirement naming**).

```text
for each Pn in order (or user subset):
  scope to Pn → TDD steps 1→8 → on GREEN auto-continue → stop when done
```

| Hard rule | |
|-----------|---|
| One phase = one full TDD cycle | No design-all-then-implement-all |
| Scope to phase exit criteria | Later work only if plan allows stubs/seams |
| Mode per phase | Classic vs backfill may differ |
| Paths / trees | **Requirement naming**; normal `./tests/<feature>/` (no phase subdirs) |
| Seal once per cycle | No unjustified rewrite of prior sealed asserts |
| Auto-continue | Then short multi-phase summary |

# Delegating to roles

Spawn/resume via runner task tool; pass requirement path — not the role prompt.
First step each sub-agent: designer `doctest skill designer --show`;
implementer `doctest skill implementer --show`. Prefix `[designer]` /
`[implementer]` when supported. Resume **same** session per role. Wait
patiently (≥1h timeout if required).

# Workflow (8 TDD steps)

## TDD step 1 — Requirements

Brainstorm both modes; write requirement; get approval. Mode per **Modes**.
Plan phase active → scope `Pn` (**Plan phases**). Tell user: (1) models/layout
if any (2) scenarios + expected output (3) layer map (**Test layers**)
(4) classic vs backfill + why (5) `Pn` or single-cycle. CLI stdout: trailing
`\n` after last content line; assert templates: newline before closing
backtick. Bugs: explore then designer; classic → RED doctests; backfill (fix
applied) → GREEN OK.

## TDD step 2 — Delegate Test Design

Write design file (**Requirement naming**); spawn designer (**Delegating**);
wait for `./tests/<feature>/`.

## TDD step 3 — Designer Questions (optional)

Resume same designer; escalate domain to user; until complete.

## TDD step 4 — Vet then run

```sh
doctest vet ./tests/<feature>
doctest test ./tests/<feature>
```

`vet` must pass (`__DOCTEST_VERSION__` — SPEC below). Interpret test per
**Modes**: classic → RED (any GREEN → re-examine); backfill → GREEN expected
(re-examine vacuous/wrong only); mixed → seal as-is.

## TDD step 5 — Seal (once)

```
git add ./tests/<feature>
```

Tests only, once per cycle. No git repo → ask before unsealed. Mixed seal as-is.

## TDD step 6 — Implement

Backfill + all sealed GREEN → skip to step 8. Else write implement file
(**Requirement naming**); spawn implementer (**Delegating**); wait GREEN. Do
not weaken already-GREEN sealed asserts.

## TDD step 7 — Implementer Questions (optional)

Resume same implementer until pass.

## TDD step 8 — Verify

```sh
git diff ./tests/<feature>            # must be clean
doctest vet ./tests/<feature>
doctest test ./tests/<feature>/...    # must be GREEN
doctest test --label "ui-automation" ./tests/<feature>/... # if ASSERT has label
doctest test --label 'slow && ui-automation' ./tests/<feature>/... # &&, ||, ()
doctest test ./...                    # no regressions
```

RED → resume implementer. Test-file edits only with explicit justification.
Report count + accepted mods. More plan phases → **Plan phases**; else summary.

# Requirement naming

Under `/tmp/` (not repo root):

| Role | Single-cycle | Plan phase N |
|------|--------------|--------------|
| Design | `/tmp/REQUIREMENT-DESIGN-<slug>.md` | `/tmp/REQUIREMENT-DESIGN-PHASE-{N}-<slug>.md` |
| Implement | `/tmp/REQUIREMENT-IMPLEMENT-<slug>.md` | `/tmp/REQUIREMENT-IMPLEMENT-PHASE-{N}-<slug>.md` |

`PHASE-{N}` only when multi-phase — never invent for single-cycle.

**Must include** (not restated in steps): mode (+ **Backfill handoff** if
backfill); layer map (**Test layers**); apply **Test-first code** +
**Parallel-safe suite** (pointers); multi-phase: `Pn` goal/exit/out-of-scope;
implement: context, tree, **"tests are sealed — do not modify"**, verify cmd,
GREEN vs RED leaves.

# Followup

Restart at TDD step 1 (relevant plan phase if multi-phase). Separate designer /
implementer sessions; reuse each role’s session ID within that role.

__DOCTEST_SPEC__

--end of skill doctest-tdd--
