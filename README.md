> I once thought every domain of test needs a dedicated DSL, now they just become markdowns with code annotation

# doctest

Doc-style test tool: write test cases as markdown decision trees with embedded Go code, then build and run them.

# Installation

```sh
go install github.com/xhd2015/doctest/cmd/doctest@latest
```

# Usage

Paste the follow prompt into your agent(claude code, codex, opencode etc.):

```md
Follow the guideline of `doctest skill doc-spec show` and `doctest skill code-spec show`, write tests first, run doctest test to ensure them RED(TDD), then seal them(`git add <tests>`); then implement the feature and run doctest until all GREEN.

<your feature here>
```


```sh
doctest test -v ./...
```

# Install skills
```sh
doctest skill --list

# the primary tdd-flow skill you need
doctest skill tdd install --codex --opencode
```

## Commands

```
Usage: doctest <command> [options]

Commands:
  agent generate <idea> [-d|--dir <target-dir>] [--agent-runner RUNNER]
  agent fill-code <target-dir>
  agent implement <prompt> [--agent-runner RUNNER] [--mock-config PATH]
  validate <dir>
  build <dir>
  test <dir>
  skill --list
  skill doc-spec show|install
  skill code-spec show|install
```

### validate

Check a directory tree conforms to the doc-style test specification
(every directory has SETUP.md, leaves have ASSERT.md, etc.).

```sh
doctest validate tests/my-feature
```

### build

Validate embedded Go code compiles (no execution). Supports the same
`<dir>` patterns as `test`, including `./...` and `./<prefix>/...`.

```sh
doctest build [-v|--verbose] [--rm] [--gen-dir DIR] [-count=N] <dir>
```

### test

Build and run all executable leaves in a doc-style test tree.

```sh
doctest test [-v|--verbose] [--rm] [-count=N] [<dir> | ./... | ./<dir>/...]
```

Examples:
```sh
doctest test -v ./
doctest test -v ./path

# Run all doctest trees under the current module
doctest test -v ./...

# Run only tests under tests/feature-a/
doctest test -v ./tests/feature-a/...
```

With `./...` or `./sub-path/...`, `doctest` walks subdirectories to find modules:

```
cwd/              (no go.mod)
├── .gitignore    # ign_a/
├── ign_a/        (gitignored → skipped)
├── mod_a/        (go.mod + DOCTest → found)
└── group/        (no go.mod)
    ├── pkg1/     (go.mod + DOCTest → found)
    └── pkg2/     (go.mod + DOCTest → found)
```

### agent generate

Use an LLM agent to generate a doc-style test tree from a feature idea.

```sh
doctest agent generate "a CLI tool that validates JSON files" -d tests/my-feature
```

### agent fill-code

Add executable Go code blocks to all SETUP.md and ASSERT.md files
in a test tree (replacing stubs with real code).

```sh
doctest agent fill-code tests/my-feature
```

### agent implement

Spawn a sub-agent to implement feature code that makes existing doc-style
tests pass. Follows the adversarial TDD workflow (main agent writes and
seals tests, sub-agent implements).

```sh
doctest agent implement "Implement JSON validation for the test tree at tests/my-feature"
```

The sub-agent can yield questions via `yield-pending-questions`; the caller
answers and re-invokes to continue the session.

### skill

Show or install specification documents as IDE skills.

```sh
doctest skill --list
doctest skill doc-spec show|install
doctest skill code-spec show|install
```

## Architecture

```
doctest/
├── main.go                    # entry point, yield-pending-questions dispatch
├── doc/                       # embedded spec documents (doc-style-test, code-spec)
│   ├── doc.go                 # //go:embed + Content() accessor
│   └── doc_test.go
├── libdoc/
│   ├── cli/                   # CLI argument parsing and dispatch
│   │   ├── cli.go
│   │   └── cli_test.go
│   ├── agent/                 # agent generate and fill-code logic
│   │   ├── agent.go
│   │   └── agent_test.go
│   ├── implementer/           # agent implement + yield-pending-questions protocol
│   │   └── implement.go
│   ├── runner/                # build/test runner (delegates to libdoc/build and libdoc/core)
│   │   ├── runner.go
│   │   ├── runner_test.go
│   │   ├── resolve.go
│   │   └── resolve_test.go
│   ├── validate/              # tree structure validation
│   │   ├── validate.go
│   │   └── validate_test.go
│   └── spec/                  # skill document embed and install
│       ├── spec.go
│       └── spec_test.go
├── tests/                     # doc-style integration tests (self-hosting)
│   ├── DOCTEST.md
│   ├── SETUP.md
│   ├── help/
│   ├── build/
│   ├── test/
│   ├── validate/
│   ├── agent/
│   ├── skill/
│   ├── main-orchestrator/
│   └── ...
└── README.md
```

### Key Design Decisions

- **Unified binary** — `yield-pending-questions` is dispatched via `os.Args[0]`, no separate build step.
- **Doc-style tests** — test cases are markdown decision trees, not Go test files. The `build`/`test` commands compile them into Go test binaries.
- **Adversarial TDD** — main agent writes and seals tests; sub-agent implements code. The main agent verifies test integrity on completion.
- **Agent implement sessions** — stored in `~/.agent-pro/data/doctest-agent/<thread-id>/`. Continuity via `CODEX_THREAD_ID`.
- **Self-hosting** — the `tests/` directory contains doc-style integration tests that exercise the doctest CLI itself.

## Doc-Style Test Specification

A doc-style test is a directory tree where:

- Every directory has a `SETUP.md` describing preconditions and steps.
- Directories with assertions have an `ASSERT.md` (these are "runnable leaves").
- `SETUP.md` may embed Go code (`type Request`, `type Response`, `func Setup`, `func Run`).
- `ASSERT.md` embeds `func Assert` to verify outcomes.
- `DOCTEST.md` at the root provides an overview and run instructions.

Run `doctest skill doc-spec show` for the full specification.

## Development

```sh
# Unit tests
go test ./libdoc/...

# Doc-style integration tests (self-hosting)
go run . test -v ./tests

# Build the binary
go build .
```
