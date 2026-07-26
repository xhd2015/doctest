---
name: doc-style-test-specification
description: when user mentions doc-style test
---

A doc-style test is a test case expressed entirely as markdown documents
organized in a decision tree. Each document describes **what** to test in
prose; the tree structure models the decision flow from feature to behavior
to scenario.

**Where a case belongs** (before growing a tree): pure flat edges → go test (L1);
multi-factor / short CLI scenarios → **in-process** doctest (L2); full process-boundary
integration only → sparse e2e (L3, `label: e2e`). See
`doctest skill design-principle --show`. Executable harness rules:
`doctest skill code-spec --show`.

## Directory Layout

```
<feature>-test-cases/
├── DOCTEST.md          # Overview: graph diagram + text tree + test case index
│                       #   Must include a header section with how to run the tests
├── SETUP.md            # Root setup — shared state and preconditions for all tests
├── testdata/           # Fixtures: minimal tree structures for validation testing
│   └── <fixture>/
│       ├── SETUP.md
│       ├── leaf/
│       │   ├── SETUP.md
│       │   └── ASSERT.md
│       └── ...
├── verify-something/   # Runnable leaf: one concrete test scenario
│   ├── SETUP.md        # Leaf-specific preconditions and steps
│   └── ASSERT.md       # Expected outcomes for this leaf only
├── decision-node/      # Abstract grouping node (no ASSERT.md)
│   ├── SETUP.md        # Inherits ancestor setup + adds node-specific setup
│   └── leaf-a/
│       ├── SETUP.md
│       └── ASSERT.md
└── ...
```

- Any directory containing an `ASSERT.md` is a **runnable leaf**.
- Directories without `ASSERT.md` are **abstract grouping/decision nodes**, and must have a `SETUP.md`
- Every runnable leaf **must** have its own `SETUP.md`.

### Why Every Directory Needs a SETUP.md

The entire design is about **drilling down** from a big, general case into
specific, concrete scenarios. Each `SETUP.md` along the path narrows the
scope by adding explicit conditions:

```
Root SETUP.md
  → "A git repository exists at /tmp/test"

  mode-commit/SETUP.md
    → "We are testing the commit command (not push, not merge, not rebase)"

    decision-untracked/SETUP.md
      → "The file being committed is untracked (not modified, not staged)"

      leaf/SETUP.md
        → "The file path contains spaces and unicode characters"

      leaf/ASSERT.md
        → "Assert: commit succeeds and the file appears in git log"
```

Without a `SETUP.md` at a given level, that directory contributes **nothing**
to distinguish its subtree from its siblings. A directory with neither
`SETUP.md` nor `ASSERT.md` is meaningless — it adds no conditions, contains no
assertions, and serves no purpose in the tree.

In short: **every directory introduces a condition that makes its subtree
different from every other branch.** That condition is recorded in its
`SETUP.md`. Leaves add `ASSERT.md` to assert the outcome of the specific
chain of conditions that leads to them.

### The testdata/ Exception

The `testdata/` directory is skipped by tree discovery. It holds fixture
directories used by other test cases, not test cases themselves.

## Decision Tree Construction

The goal is to exhaustively cover every behavior, every input variant,
every edge case, and every error path of the feature under test.

### 1. Analyze the Input

Examine the feature from every angle. For a CLI command, consider:

- **Every flag** — presence, absence, and every valid/invalid value
- **Every positional argument** — present, absent, valid, invalid, edge case values
- **Every combination** — flags that interact, args that depend on flags, mutually exclusive options
- **The environment** — OS state, file system, network conditions, dependencies

For an API, consider every parameter, every HTTP method, every status code.
For a data structure, consider every field, every type, every valid/invalid range.

### 2. Exhaustive Decomposition

Break the feature into a tree of decisions:

- **Modes** — top-level branches: subcommands, operation types, major flag groupings.
  Each mode is a first-level directory under root.
- **Decision points** — conditions that fork behavior inside a mode: boolean flags,
  type checks, state queries, input validation outcomes. Each decision is a
  subdirectory under its mode.
- **Error states** — every path that produces an error: invalid input, missing
  prerequisites, conflict conditions, unexpected failures.
- **Happy-path outcomes** — the expected success paths and their observable results.

Prefer **more leaves over fewer**. When in doubt, split into separate test cases.
Each leaf tests exactly one scenario.

### 3. Propose and Get Approval

Before building the full tree, present a flat list of proposed test cases
(name + one-line description). Get explicit approval before writing files.

### 4. Write the Tree

Create the directory structure. Write each document with verbose, explicit prose.

### Root SETUP.md: Target Behavior Overview

The root `SETUP.md` serves a dual role. In addition to shared preconditions,
it documents **how the target program or function behaves** — a specification of
the thing under test:

- What the tool does (its purpose, its inputs, its outputs)
- Every flag, option, or parameter and what each controls
- Every subcommand or operation mode and how they differ
- Expected outcomes for common scenarios (happy path)
- Error conditions and how they are surfaced (error messages, exit codes)

This overview sets the stage for the entire tree. Child and leaf `SETUP.md`
files then drill down into specific branches of this behavior:

```
Root SETUP:          "commit creates a snapshot of staged changes with a message"
  mode-message/:     "commit requires a non-empty message"
    leaf-empty/:     "an empty message is rejected with exit code 1"
    leaf-whitespace/:"a whitespace-only message is rejected"
  mode-paths/:       "commit handles file paths of various forms"
    leaf-unicode/:   "unicode filenames are committed correctly"
```

The root describes the surface area; the children decompose it into testable
units.

## SETUP.md Format

SETUP.md describes the preconditions and steps needed to reach the test state.
It has structured sections:

```markdown
## Preconditions
- Explicit state that must exist before this test runs
- e.g. "a git repository exists at /tmp/test with two commits on main"
- e.g. "the user 'admin' with password 'secret' is registered"

## Steps
1. First action the tester performs
2. Second action
3. ... (implicitly inherited from ancestor SETUP.md files)

## Context
- Environment details: OS, tools, configuration
- Roles and permissions
- External service state
```

### Verbose and Explicit

Every description must be **concrete and reproducible**. Never write:
- "test the login flow"
- "try with bad input"
- "check it works"

Instead write:
- "Given a registered user with email `alice@example.com` and password `s3cret!`, when login is called with those credentials, then the response contains a valid JWT with subject `alice@example.com`"
- "Given the `--force` flag is NOT set, when move is called with a target path that already exists, then the command exits with code 1 and stderr contains 'target already exists'"

### Inheritance Rules

Two distinct rules govern how documents interact across the tree:

#### SETUP Chains

`SETUP.md` accumulates along the ancestor path. A leaf's effective setup is the
**union** of all `SETUP.md` files from root to leaf:

```
Root SETUP:        "A git repo exists at /tmp/test"
  mode-commit/:    "We are testing the commit command"
    decision/:     "The file is untracked"
      leaf/:       "The file path contains spaces"
```

A test at `mode-commit/decision/leaf/` inherits all four levels: the git repo
precondition (root), the commit mode (mode-commit), the untracked decision
(decision), and the unicode path detail (leaf). Nothing is repeated.

- `## Steps` — concatenated root-first to leaf-last
- `## Preconditions` and `## Context` — unioned from all ancestors plus the leaf

#### ASSERT Is Case-Private

`ASSERT.md` belongs to **one leaf only**. It is never inherited, never shared,
never merged. Each leaf's assertions are self-contained and must stand on their
own. A leaf directory without `ASSERT.md` is a grouping/decision node, not a
runnable test case.

## ASSERT.md Format

ASSERT.md describes the expected outcomes for a **single** leaf.

### Optional YAML frontmatter (run profile labels)

Expensive, flaky, manual, or UI leaves may prefix the file with frontmatter:

```yaml
---
label: slow, ui-automation
explanation: needs display server; ~30s per run
---
```

- **`label`** is a **scalar YAML string**, not a YAML array. Multiple tags on one leaf
  use a **comma-separated** list (`slow, ui-automation`). Any non-empty `label` skips the
  leaf in discovery mode (`doctest test` on tree root, grouping dir, or `./...`).
- **`explanation`** is optional prose; it does **not** skip by itself.
- Canonical run-profile labels: `slow`, `heavy`, `flaky`, `manual`, `ui-automation`.

To run labeled leaves selectively, use `doctest test --label EXPR` (`&&`, `||`,
parentheses; repeatable `--label` flags are OR'd). Without `--label`, discovery runs
only unlabeled leaves; point at a concrete leaf path to run one labeled test.

### Assertion sections

It has structured sections; use only those relevant to the specific scenario:

```markdown
## Expected Output
Optional but **recommended** when asserting CLI or text output with the
`github.com/xhd2015/doctest/assert` DSL. Fenced block mirrors the template
passed to `assert.Output` — not required by `doctest vet`. See
`doctest skill output-assert --show`.

```
---
version: 3
---
Usage: mytool
  build
  test
```

## Expected
- Observable outcome that confirms success
- e.g. "stdout matches output template (see Expected Output)"
- e.g. "the response body has field 'status' equal to 'ok'"

## Side Effects
- State changes that must have occurred
- e.g. "file /a no longer exists"
- e.g. "file /b exists and its contents match /a"

## Errors
- Expected error message if this is an error test case
- e.g. "stderr contains 'permission denied'"

## Exit Code
- Expected exit code for CLI tools
- e.g. exit code 0 on success, 1 on error
```

## DOCTEST.md Format

The root `DOCTEST.md` serves as the entry point for the test suite. It must
include a header section that tells readers how to execute all tests:

```markdown
## How to Run

Document default discovery and how to run labeled leaves when the tree has any
`label:` frontmatter:

```sh
doctest test ./

# Run only slow-labeled leaves (unlabeled leaves skipped)
doctest test ./ --label slow

# Boolean expressions and multiple flags
doctest test ./ --label 'slow && ui-automation'
doctest test ./ --label slow --label heavy
```

## Verification

After creating the tree and writing all documents, verify they match the
specification by running:

```sh
validate_test_case_tree <output-dir>
```

If the tree includes no executable code, the doc structure alone is validated.
