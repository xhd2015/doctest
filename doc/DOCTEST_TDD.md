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

How you delegate to those roles depends on what your agent runner provides (see
**Choosing a delegation channel** below). The roles and the
`designer → RED → seal → implementer` flow are fixed; only the transport
changes.

If you've already reproduced a bug with written doctests, and confirm they're RED, then you can skip the design → RED stages, start directly from seal step.

# Choosing a delegation channel

Use the **first** option that applies:

| Priority | When | How |
|----------|------|-----|
| **1 — Native subagent/task** | Your runner exposes a subagent or task tool (e.g. Task, runner-native subagent API) | Spawn/resume subagents with the designer or implementer role prompt and the requirement as the task. See **Role prompts** below. |
| **2 — CLI fallback** | No native subagent mechanism, or you need doctest session hooks (`yield-pending-questions`, `--status`, progress reporting) | `doctest agent design` / `doctest agent implement` |

**Role prompts** — load once per role, then pass as the subagent system/role context:

Run via `bash` to load the role setup: **`doctest skill designer show`** and **`doctest skill implementer show`**.

**Session continuity** — designer and implementer each have their own session.
Within a role, resume the same session for follow-up questions (native: runner
resume/session ID; CLI: `--session-id`).

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

For bugs: reproduce is critical. Use an explore/analysis sub-agent (native
task/subagent when available) to narrow scope first, then delegate to the
**designer** role to write failing doctests.

If there is anything needs clarification, list them and ask for user confirmation until no gap understanding user's intent.

## Phase 2 — Delegate Test Design

Write `REQUIREMENT-DESIGN-<context-summary-and-feature-slug>.md` from Phase 1.

### Preferred — native subagent/task

Spawn a subagent with:

- **Role/system**: designer role prompt (see **Role prompts** above)
- **Task**: contents of `REQUIREMENT-DESIGN-<slug>.md` (or the file path plus a short summary)

Wait until the designer reports the doctest tree is written under
`./tests/<feature>/`.

### Fallback — `doctest agent design`

```sh
# please wait enough for any sub-agent, if timeout required, set to 1h
doctest agent design --timeout 1h --requirement REQUIREMENT-DESIGN-<context-summary-and-feature-slug>.md

# or for short requirement or followup
doctest agent design --timeout 1h <<EOF
<design doc from Phase 1>
EOF
```

## Phase 3 — Designer Questions (optional)

If the designer yields questions, answer them and **resume the same designer
session** — do not start a fresh designer.

### Preferred — native subagent/task

Resume the designer subagent with the answers (and escalate domain-specific
questions to the user first).

### Fallback — `doctest agent design`

```sh
doctest agent design --timeout 1h <<EOF
<answers to questions>
EOF
```

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

### Preferred — native subagent/task

Spawn a subagent with:

- **Role/system**: implementer role prompt (see **Role prompts** above)
- **Task**: contents of `REQUIREMENT-IMPLEMENT-<slug>.md`

Wait until the implementer reports all tests passing.

### Fallback — `doctest agent implement`

```sh
# please wait enough for any sub-agent, if timeout required, set to 1h
doctest agent implement --timeout 1h --requirement REQUIREMENT-IMPLEMENT-<slug>.md
```

To check status (CLI only):

```
doctest agent implement --session-id <ID> --status
```

## Phase 7 — Implementer Questions (optional)

Same pattern as Phase 3: resume the **same implementer session** with answers
until all tests pass.

### Preferred — native subagent/task

Resume the implementer subagent with answers or failure output from Phase 8.

### Fallback — `doctest agent implement`

Re-invoke with follow-up prompt or heredoc (same `--session-id` for continuity).

## Phase 8 — Verify

```sh
git diff ./tests/<feature>            # must be clean
doctest vet ./tests/<feature>         # structure still valid
doctest test ./tests/<feature>/...    # must be GREEN

doctest test --label "ui-automation" ./tests/<feature>/... # if the ASSERT.md contains label header
doctest test --label 'slow && ui-automation' ./tests/<feature>/... # --label accepts simple expr lie expr&&, || , ()

doctest test ./...                    # no regressions
```

If RED, feed failures back to the implementer (resume native subagent if allowed, or
re-invoke CLI). If test files were modified, accept only with explicit
justification (the test expected wrong behavior per spec).

Report: test count, modifications accepted (with rationale).

# Requirement File Naming

- Design: `REQUIREMENT-DESIGN-<slug>.md`
- Implement: `REQUIREMENT-IMPLEMENT-<slug>.md`

Pass the short prompt or long requirement file path (plus optional summary) to native subagents. For the CLI
fallback, use `--requirement` for long prompts or prompts with shell-special
characters (`$`, `#`, `!`), or a heredoc for adhoc followup.

# Followup Requests

Every followup (new feature, fix, amendment) restarts this workflow from
Phase 1. Design and implement sessions are isolated — use their respective
session IDs within each role.

Designer and implementer sessions are always separate (different tasks). Once
you establish a session for a role, keep that same session ID for all
follow-ups within that role (native resume or CLI `--session-id`).

__DOCTEST_SPEC__

--end of skill doctest-tdd--