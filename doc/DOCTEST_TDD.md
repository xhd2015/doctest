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

**You NEVER touch source files.** No Edit/Write on source or config. All code
changes go through:

- **Designer** — writes doctest trees
- **Implementer** — writes implementation code

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

- Already-reproduced bug with RED doctests → skip design/RED; start at seal.
- Backfill: RED not required for leaves documenting correct behavior; GREEN
  expected. Skip implementer when all sealed tests are already GREEN.

# Delegating to roles

Spawn/resume role subagents via the runner's task tool. Pass a short requirement
(or file path) — not the role prompt.

Each sub-agent **must** run as its first step:

- Designer: `doctest skill designer show`
- Implementer: `doctest skill implementer show`

Prefix descriptions with `[designer]` / `[implementer]` when supported. Resume
the **same** session per role for follow-ups. Wait patiently (use ≥1h timeout
if the runner requires one).

# Workflow (8 Phases)

## Phase 1 — Requirements

Brainstorm (both modes). Produce a requirement file; get explicit approval.

Auto-detect mode (see **Modes**). State it in the requirement file.

Tell the user:

1. Data models and storage layout (if any)
2. Scenarios and expected output
3. How you will test (prefer doctests)
4. **Classic TDD** or **coverage backfill**, and why

CLI: user-facing stdout ends with `\n` after the last content line; when using
doctest assert templates, newline before the raw string's closing backtick.

Bugs: explore/narrow scope first, then designer. Classic → failing doctests;
backfill (fix applied) → doctests of fixed behavior (GREEN OK).

Clarify until intent is clear.

## Phase 2 — Delegate Test Design

Write `REQUIREMENT-DESIGN-<slug>.md`. Spawn designer with that path (+ optional
summary). Designer runs `doctest skill designer show` first.

**Backfill — MUST tell the designer** (spawn message and/or requirement):

- Mode: coverage backfill (implementation present)
- Intent: **backfill** missing doctests for existing correct behavior
- RED not required; GREEN expected for covered paths
- Mixed GREEN/RED OK when some behaviors still missing
- Do not invent must-fail assertions only for classic TDD theater

Wait until the tree is under `./tests/<feature>/`.

## Phase 3 — Designer Questions (optional)

Answer questions by **resuming the same designer session**. Escalate
domain questions to the user first. Repeat until designer completes.

## Phase 4 — Vet then run

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

## Phase 5 — Seal (once)

```
git add ./tests/<feature>
```

Seal tests only, once. Outside a git repo, ask before proceeding unsealed.
**Mixed suites seal as-is** (GREEN + RED together).

## Phase 6 — Implement

If backfill and all sealed tests are GREEN → **skip** (no implementer); go to
Phase 8.

Else write `REQUIREMENT-IMPLEMENT-<slug>.md`: context, summary, tree structure,
**"tests are sealed — do not modify"**, verify command, which leaves were
already GREEN vs still RED.

Spawn implementer with that path. Implementer runs `doctest skill implementer
show` first. Wait until all tests pass. Do not weaken already-GREEN sealed
asserts.

## Phase 7 — Implementer Questions (optional)

Resume the **same implementer session** with answers or Phase 8 failures until
all pass.

## Phase 8 — Verify

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

# Requirement File Naming

- Design: `REQUIREMENT-DESIGN-<slug>.md` — must state mode; backfill includes
  Phase 2 handoff bullets
- Implement: `REQUIREMENT-IMPLEMENT-<slug>.md`

# Followup Requests

Every followup restarts at Phase 1. Keep designer and implementer sessions
separate; reuse each role's session ID within that role.

__DOCTEST_SPEC__

--end of skill doctest-tdd--
