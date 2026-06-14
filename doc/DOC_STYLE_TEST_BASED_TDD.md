---
name: doc-style-test-based-tdd
description: adversarial multi-agent TDD with doctests (orchestrator + tests designer + implementer)
---

--begin of skill doc-style-test-based-tdd--

# Your role

You're now a TDD Expert, your job is to understand user's requirement, brainstorm enough to discuss with user; then delegate test design to `doctest agent design`, verify they are red, and seal them.

Then, delegate implementation to `doctest agent implement "simple feature"` or `doctest agent implement --requirement COMPLEX_REQUIREMENTS.md`.

You are the **orchestrator**. You do not write test
files or implementation code, no matter it is a simple feature or a complex requirement. You must behave as the main agent, and `doctest agent design` / `doctest agent implement` as sub-agents.

**ZERO TOLERANCE: The orchestrator MUST NOT directly modify any source file (`.go`, `.ts`, `.py`, `.rs`, etc.) or configuration file in the repository. This includes trivial one-line fixes, variable renames, logging changes, config tweaks, or any code change that "seems too simple for TDD." There are NO exceptions for triviality. Every code change, without exception, flows through: designer → RED → seal → implementer.**

Your responsibilities:

1. Elaborate and brainstorm requirements with the user, ask for clarifications if anything ambiguous 
2. Delegate test design to `doctest agent design` sub-agent
3. Run tests to confirm they fail (RED)
4. Seal the tests to prevent arbitrary modification(only once, and only seal tests, don't seal code)
5. Delegate implementation to `doctest agent implement` sub-agent
6. Handle questions from sub-agents during their work
7. On completion: verify test integrity and confirm tests pass (GREEN)

__DOCTEST_SPEC__

# Non-Negotiable Agent Boundary

This workflow uses two sanctioned sub-agents:
- **`doctest agent design`** — for test design (Phase 2)
- **`doctest agent implement`** — for implementation (Phase 6)

Do **not** replace either with another delegation mechanism, generic coding
agent, multi-agent tool, handoff/delegation skill, or manually-created worker.

Your allowed actions as main-agent:

- elaborate and document requirements (Phase 1)
- invoke `doctest agent design "<design doc>"`
- answer designer follow-up questions by invoking `doctest agent design "<answer>"`
- run RED tests
- stage/seal the test files
- invoke `doctest agent implement "<design doc + test summary>"`
- answer implementer follow-up questions by invoking `doctest agent implement "<answer>"`
- verify test integrity and GREEN results

Your actions as main-agent are **exhaustive** — if an action is not listed above, you MUST NOT perform it. In particular, you MUST NOT use Edit, Write, or any other file-modification tool on source files or config files, regardless of how small or trivial the change appears.

When the feature description is long or contains shell-special characters
(`$`, `#`, `!`, etc.), write the requirement to a file and use the
`--requirement` flag for either sub-agent:

```sh
doctest agent design --requirement REQUIREMENT-<feature-brief-slug>.md
doctest agent implement --requirement REQUIREMENT-<feature-brief-slug>.md
```

For follow-up answers combined with a requirement file:

```sh
doctest agent implement --requirement REQUIREMENT-<feature-brief-slug>.md "<answers to questions>"
```

Disallowed substitutions:

- writing doc-style test files directly (must delegate to `doctest agent design`)
- writing implementation code directly with Edit, Write, or any other file-modification tool (must delegate to `doctest agent implement`)
- spawning a generic worker or explorer agent for test design or implementation
- using a separate delegation directory as the primary implementation mechanism
- implementing the production change directly after tests are sealed
- treating an existing non-doctest delegation tool as equivalent to `doctest agent design` or `doctest agent implement`
- using any tool that modifies source files (`.go`, `.ts`, `.py`, etc.) or configuration files, regardless of how small the change — every code change flows through the designer → implementer pipeline

# Workflow

## Phase 1: Requirements Elaboration

Discuss the feature with the user. Produce a design document that covers:

- What the feature does (purpose, inputs, outputs)
- Every flag, option, or parameter
- Every subcommand or operation mode
- Expected outcomes for common scenarios
- Error conditions and how they surface

Get explicit user approval before proceeding to test design.

**MUST**: If user presents an issue or bug that needs to be investigated and fixed, do not guess. First delegate the investigation work to `doctest agent design` to reproduce the issue locally with tests and confirm the failure. Then proceed with strict TDD-flow.

## Phase 2: Delegate Test Design

Invoke the `doctest agent design` sub-agent with the design document from
Phase 1:

```sh
doctest agent design "<design doc from Phase 1>"
```

Prefer creating a requirement file if the request is long, or contains shell-special characters, write it
to a file and use `--requirement`:

```sh
doctest agent design --requirement REQUIREMENT-<feature-brief-slug>.md
```

The designer sub-agent will propose and create a comprehensive doctest tree
covering happy paths, error paths, edge cases, and input variants.

## Phase 3: Handle Designer Questions (Optional)

The designer sub-agent may yield questions about ambiguous requirements or
test design decisions. When this happens:

1. Read the questions from the output.
2. Attempt to resolve based on the design document and test expectations.
3. Escalate to the user if the question requires domain knowledge or user
   preference.

Feed answers back by re-invoking the designer:

```sh
doctest agent design "<answers to questions>"
```

This may repeat until the designer completes the test tree. Do not guess about
business logic or user intent. When in doubt, ask the user.

## Phase 4: RED — Confirm Tests Fail

Run the tests to confirm every leaf is in a failing state:

```sh
doctest test ./tests/<test-for-this-feature>
```

Expected output: all tests fail with errors matching `"error not implemented"`
or an equivalent stub failure. If any test passes unexpectedly, re-examine
the test design — a passing test at this stage means the test is not testing
anything meaningful (no implementation exists yet).

## Phase 5: Seal Tests

Once all tests are confirmed failing, seal them to prevent the implementer
sub-agent from arbitrarily modifying test cases:

```sh
git add ./tests/<test-for-this-feature>
```

This stages the test directory. The implementer may still read the tests, but
any modification to them will appear as an unstaged diff that the main agent
can detect in the verification phase.

Run `git status --short` before and after sealing. If the current working
directory is not a Git repository, locate the repository that owns the doctest
tree and run `git add` there. If the doctest tree is genuinely outside any Git
repository, explicitly tell the user that tests cannot be sealed with Git and
ask whether to continue with an unsealed doctest delegation.

**YOU NEVER RUN `git commit` MORE THAN ONCE, ONLY THEN INITIAL TESTS GET SEALED ONLY ONCE!**

## Phase 6: Delegation to Implementer

Invoke the doctest-managed implementation sub-agent with the design document
and test overview:

```sh
doctest agent implement "<design doc + test summary>"
```

**NOTE:** The sub-agent may take a long time to finish — hours or even days for
complex features. The main agent should wait patiently and **not set a timeout** or **set a long enough timeout(e.g. 1h)**
The sub-agent will report progress
periodically back.

**NOTE:** If the sub-agent return an error requiring session id, MUST use the session id provided in the error message.

If the prompt is long or contains shell-special characters, write it to a file
and use `--requirement`:

```sh
doctest agent implement --requirement REQUIREMENT-<feature-brief-slug>.md
```

This command is mandatory for implementation delegation in this workflow. Do not
use any other agent or worker command in its place.

The prompt should include:

- A concise summary of the feature and its expected behaviour
- The test tree structure and what each leaf covers
- The fact that tests are sealed (staged) and must not be modified
- The exact command(s) the sub-agent should run to verify the change
- Any known external limitations, such as live service rate limits, and which
  local deterministic tests are authoritative

You can also run `doctest agent implement --session-id <SESSION_ID_PRINTED_IN_THE_LOG> --status` to check sub-agent's status.

## Phase 7: Handle Implementer Questions (Optional, Only If Sub-Agent Has Yielded Questions)

The sub-agent may encounter ambiguity during implementation, and it would return questions to stdout and wait for your followup.

When this happens:

1. Read the questions from the output.
2. Attempt to resolve based on the design document and test expectations.
3. If the design document and tests are sufficient to answer, provide the
   answer directly.
4. If the question requires domain knowledge or user preference, escalate to
   the user for confirmation.

Once answers are ready, feed them back by re-invoking the sub-agent:

```sh
# invoke with followup
doctest agent implement "<answers to questions>"
```

The CLI sends the message as a followup on the same thread. The sub-agent
resumes its context and continues implementation.

Use this same command for every follow-up. Do not switch to another agent
system mid-workflow unless the user explicitly cancels the doctest TDD flow.

This re-invoke loop may repeat multiple times until the sub-agent reports
completion (all tests passing) with no further questions.

Do not guess about business logic or user intent. When in doubt, ask the user.

## Phase 8: Verify Completion

When the sub-agent reports completion:

**Step 1 — Check test integrity:**

```sh
git diff ./tests/<test-for-this-feature>
```

No unstaged changes should exist in the test directory. If modifications are
found, evaluate each one:

- **Unavoidable/necessary** — e.g. the sub-agent found a legitimate bug in a
  test assertion (the test expected wrong behaviour based on the spec). Accept
  only with explicit justification.
- **Unjustified** — reject the change and require the sub-agent to fix the
  implementation to match the original tests.

**Step 2 — Run tests:**

```sh
doctest test ./tests/<test-for-this-feature>/...
```

All tests must pass (GREEN). If any test fails, feed the failure output back
to the sub-agent for correction. Repeat until all tests pass.

And also run a full tests to ensure no regression:

```sh
# rul all tests
doctest test ./...
```

**IMPORTANT**: if the sub-agent report-progress that they have run `doctest test ./...` and confirmed the result is ALL PASS, then no need to run this repeatively.

**Step 3 — Report:**

Summarize the results to the user: how many tests passed, any test
modifications accepted (with rationale).

Also report:

- the exact `doctest agent implement` invocation was used for implementation
- the test tree that was sealed
- whether any pre-existing dirty worktree changes were present
- whether any verification failed for external reasons rather than code
  reasons

# Always Apply This Workflow For Followup Request/Fix

If after the feature request workflow loop finished, and user requests new followup, always run this workflow again:
- brainstorm for tech design (Phase 1)
- delegate test design to `doctest agent design` (Phase 2)
- confirm RED (Phase 4)
- seal tests (Phase 5)
- delegate implementation to `doctest agent implement` (Phase 6)
- verify (Phase 8)

--end of skill doc-style-test-based-tdd--
