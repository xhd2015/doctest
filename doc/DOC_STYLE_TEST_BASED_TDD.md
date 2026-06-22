---
name: doc-style-test-based-tdd
description: adversarial multi-agent TDD with doctests (orchestrator + tests designer + implementer)
---

--begin of skill doc-style-test-based-tdd--

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

All code changes happen exclusively through two sub-agents:
- **`doctest agent design`** — writes tests
- **`doctest agent implement`** — writes implementation

Do not replace them with generic workers, handoff agents, or any other
delegation mechanism.

# Workflow (8 Phases)

Every feature request, bug fix, or followup follows this loop.

## Phase 1 — Requirements

Brainstorm with user. Produce a requirement file. Get explicit user approval
before continuing. 

Explicitly tell user:
1. What the underlying data models and storage layout(if any) are;
2. What scenarios you will test, and expected output;
3. How you gonna test that, prefer rerunable tests(doc-style tests or unit tests);

For bugs: reproduce is critical; use an explore/analysis sub-agent to narrow scope
first, then delegate to `doctest agent design` to reproduce.

If there is anything needs clarification, list them and ask for user confirmation until no gap understanding user's intent.

## Phase 2 — Delegate Test Design

```sh
# please wait enough for any sub-agent, if timeout required, set to 1h
doctest agent design --timeout 1h --requirement REQUIREMENT-DESIGN-<context-summary-and-feature-slug>.md

# or for short requirement or followup
doctest agent design --timeout 1h <<EOF
<design doc from Phase 1>
EOF
```

Wait patiently. Do not set a timeout; use ≥1h if needed.

## Phase 3 — Designer Questions (optional)

```sh
doctest agent design --timeout 1h <<EOF
<answers to questions>
EOF
```

Escalate to user for domain-specific questions. Repeat until the designer
completes.

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

```sh
# please wait enough for any sub-agent, if timeout required, set to 1h
doctest agent implement --timeout 1h --requirement REQUIREMENT-IMPLEMENT-<slug>.md
```

The requirement file must include: summarized context, feature summary, test
tree structure, **"tests are sealed — do not modify"**, and the verify
command.

To check status:
```
doctest agent implement --session-id <ID> --status
```

## Phase 7 — Implementer Questions (optional)

Same pattern as Phase 3, using `doctest agent implement`. Re-invoke until the
implementer reports all tests passing.

## Phase 8 — Verify

```sh
git diff ./tests/<feature>            # must be clean
doctest vet ./tests/<feature>         # structure still valid
doctest test ./tests/<feature>/...    # must be GREEN
doctest test ./...                    # no regressions
```

If RED, feed failures back to implementer. If test files were modified, accept
only with explicit justification (the test expected wrong behavior per spec).

Report: test count, modifications accepted (with rationale).

# Requirement File Naming

- Design: `REQUIREMENT-DESIGN-<slug>.md`
- Implement: `REQUIREMENT-IMPLEMENT-<slug>.md`

Use `--requirement` for long prompts or prompts with shell-special characters
(`$`, `#`, `!`).

Use heredoc for adhoc followup.

# Followup Requests

Every followup (new feature, fix, amendment) restarts this workflow from
Phase 1. Design and implement sessions are isolated — use their respective
session IDs.

For `doctest agent design/implement`, their sessions are isolated(because they do different tasks), so their session ids are different. And if you've passed session id once, for followup, please also keep that same session id for continuity.

Always wait subagents patiently. Do not set a timeout; use ≥1h if needed. Sub-agent reports
progress periodically.

__DOCTEST_SPEC__

--end of skill doc-style-test-based-tdd--
