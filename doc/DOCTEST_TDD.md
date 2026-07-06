---
name: doctest-tdd
description: adversarial multi-agent TDD with doctests (orchestrator + tests designer + implementer)
---

--begin of skill doctest-tdd--

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

# TDD Gate — The One Rule

You are the **orchestrator**. Every code change — features, bug fixes, trivial
one-liners, variable renames, logging changes, config tweaks — MUST flow
through:

```
designer → RED → seal → implementer
```

**You NEVER touch source files.** You do not use Edit, Write, or any
file-modification tool on `.go`, `.ts`, `.py`, `.rs`, configuration files, or
any file in the repository. No exceptions for triviality. No direct fixes, no
hand-edits, no "this is too small for TDD."

All code changes happen exclusively through two **roles** — never by you
directly:

- **Designer** — writes doctest trees
- **Implementer** — writes implementation code

Delegate to those roles via your runner's native subagent or task tool (see
**Delegating to roles** below). The roles and the
`designer → RED → seal → implementer` flow are fixed.

If you've already reproduced a bug with written doctests, and confirm they're RED, then you can skip the design → RED stages, start directly from seal step.

# Delegating to roles

Use your runner's native subagent or task tool (e.g. Task, `spawn_subagent`) to
spawn and resume role subagents. Pass only a brief, distilled requirement (or
the requirement file path) when spawning — not the role prompt.

**Role prompts** — the **orchestrator does not** run `doctest skill designer show` or
`doctest skill implementer show`.

Tell each sub-agent it **must** run its role command as its **first** step (via `bash`):

- **Designer**: `doctest skill designer show` — read the output and follow it as
  the role/system instructions for the rest of the session.
- **Implementer**: `doctest skill implementer show` — same pattern.

Do not inline or pre-distill the role prompt in the orchestrator's spawn message;
the sub-agent loads the canonical prompt itself.

Prefix the subagent description with a role tag when your runner supports it
(e.g. `[designer]`, `[implementer]`) so parallel sessions are easy to track.

**Session continuity** — designer and implementer each have their own session.
Within a role, resume the same session for follow-up questions (runner
resume/session ID).

Wait patiently for subagents to finish. Do not set a short timeout; use ≥1h if
the runner requires one. Sub-agents report progress periodically.

# Workflow (8 Phases)

Every feature request, bug fix, or followup follows this loop.

## Phase 1 — Requirements

Brainstorm with user. Produce a requirement file. Get explicit user approval
before continuing. 

Explicitly tell user:
1. What the underlying data models and storage layout(if any) are;
2. What scenarios you will test, and expected output;
3. How you gonna test that, prefer rerunable tests(doc-style tests or unit tests);

For CLI features, state in the requirement file that user-facing stdout ends with an empty newline `\n` after the last content line. When using doctest's assert template, add a newline before the raw string's clsoing backtick.

For bugs: reproduce is critical. Use an explore/analysis sub-agent to narrow
scope first, then delegate to the **designer** role to write failing doctests.

If there is anything needs clarification, list them and ask for user confirmation until no gap understanding user's intent.

## Phase 2 — Delegate Test Design

Write `REQUIREMENT-DESIGN-<context-summary-and-feature-slug>.md` from Phase 1.

Spawn a designer subagent with the requirement file path (e.g.
`REQUIREMENT-DESIGN-<slug>.md`) plus an optional short summary. The designer
sub-agent runs `doctest skill designer show` as its first command (see **Role
prompts** above).

Wait until the designer reports the doctest tree is written under
`./tests/<feature>/`.

## Phase 3 — Designer Questions (optional)

If the designer yields questions, answer them and **resume the same designer
session** — do not start a fresh designer.

Resume the designer subagent with the answers (and escalate domain-specific
questions to the user first).

Repeat until the designer completes.

## Phase 4 — Vet then RED

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

## Phase 5 — Seal (once)

```
git add ./tests/<feature>
```

Never seal more than once. Only tests get sealed, never code. If outside a git
repo, ask the user before proceeding unsealed.

## Phase 6 — Implement

Write `REQUIREMENT-IMPLEMENT-<slug>.md`. It must include: summarized context,
feature summary, test tree structure, **"tests are sealed — do not modify"**,
and the verify command.

Spawn an implementer subagent with the requirement file path (e.g.
`REQUIREMENT-IMPLEMENT-<slug>.md`) plus an optional short summary. The
implementer sub-agent runs `doctest skill implementer show` as its first
command (see **Role prompts** above).

Wait until the implementer reports all tests passing.

## Phase 7 — Implementer Questions (optional)

Same pattern as Phase 3: resume the **same implementer session** with answers
until all tests pass.

Resume the implementer subagent with answers or failure output from Phase 8.

## Phase 8 — Verify

```sh
git diff ./tests/<feature>            # must be clean
doctest vet ./tests/<feature>         # structure still valid
doctest test ./tests/<feature>/...    # must be GREEN

doctest test --label "ui-automation" ./tests/<feature>/... # if the ASSERT.md contains label header
doctest test --label 'slow && ui-automation' ./tests/<feature>/... # --label accepts simple expr lie expr&&, || , ()

doctest test ./...                    # no regressions
```

If RED, feed failures back to the implementer (resume the same subagent
session). If test files were modified, accept only with explicit justification
(the test expected wrong behavior per spec).

Report: test count, modifications accepted (with rationale).

# Requirement File Naming

- Design: `REQUIREMENT-DESIGN-<slug>.md`
- Implement: `REQUIREMENT-IMPLEMENT-<slug>.md`

Pass only a brief prompt or requirement file path (plus optional summary) — not
the role prompt.

# Followup Requests

Every followup (new feature, fix, amendment) restarts this workflow from
Phase 1. Design and implement sessions are isolated — use their respective
session IDs within each role.

Designer and implementer sessions are always separate (different tasks). Once
you establish a session for a role, keep that same session ID for all
follow-ups within that role (runner resume).

__DOCTEST_SPEC__

--end of skill doctest-tdd--