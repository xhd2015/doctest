# TDD Tests: `doctest agent implement`

This document defines the test scenarios to be created **ahead of implementation**.
Each test expects failure (RED) until the implementation exists.

---

## Test Tree Location

```
agents/doctest/libdoc/agent/tests/
```

Tests use `fake-opencode` via `--mock-config` — no real LLM calls.

---

## Test Tree Structure

```
tests/
├── DOCTEST.md                          # how to run: doctest test ./
├── SETUP.md                            # root harness: build doctest, helpers
│
├── first-call/                         # first invocation behavior
│   ├── SETUP.md
│   ├── creates-session/
│   │   ├── SETUP.md                    # write mock config, invoke
│   │   └── ASSERT.md                   # CODEX_THREAD_ID set, session dir exists
│   └── no-prompt-error/
│       ├── SETUP.md                    # invoke with empty prompt
│       └── ASSERT.md                   # exits non-zero, stderr contains "requires"
│
├── completion/                         # sub-agent finishes successfully
│   ├── SETUP.md
│   └── sub-agent-completes/
│       ├── SETUP.md                    # mock returns success events
│       └── ASSERT.md                   # exit 0, stdout contains mock text
│
├── questions/                          # yield-pending-questions flow
│   ├── SETUP.md
│   ├── sub-agent-yields-questions/
│   │   ├── SETUP.md                    # mock hook calls yield-pending-questions
│   │   └── ASSERT.md                   # exit 0, questions in output
│   └── resume-after-questions/
│       ├── SETUP.md                    # set CODEX_THREAD_ID, invoke followup
│       └── ASSERT.md                   # exit 0, mock completes on same session
│
├── failure/                            # sub-agent error paths
│   ├── SETUP.md
│   └── sub-agent-exits-nonzero/
│       ├── SETUP.md                    # mock returns exit_code=3
│       └── ASSERT.md                   # exits non-zero, stderr contains error
│
└── yield-pending-questions-binary/     # the dispatched binary itself
    ├── SETUP.md
    ├── writes-to-fifo/
    │   ├── SETUP.md                    # create FIFO, invoke binary with args
    │   └── ASSERT.md                   # FIFO contains question JSON
    └── no-fifo-error/
        ├── SETUP.md                    # invoke binary without QUESTION_FIFO
        └── ASSERT.md                   # stderr contains "QUESTION_FIFO must be set"
```

---

## Root Harness

### `tests/SETUP.md`

Defines the shared model and helpers:

```go
type Request struct {
    TempDir       string
    DoctestBin    string   // path to built doctest binary
    MockConfigPath string  // path to mock config JSON
    Args          []string // arguments to pass to doctest
    Env           []string // extra env vars
}

type Response struct {
    ExitCode int
    Stdout   string
    Stderr   string
    Err      error
}
```

- `func Setup(t, req)` — builds `doctest` binary via `(cd agents/doctest && go build -o ... .)`
- `func Run(t, req)` — runs `req.DoctestBin agent implement ...` with `req.Args` and `req.Env`, captures output
- Helper: `writeMockConfig(t, req, body string)` — writes a mock config JSON file
- Helper: `assertContains(t, s, want)` / `assertExitCode(t, resp, code)` etc.

---

## Test Scenario Details

### 1. `first-call/creates-session/`

**What it tests**: The first call to `doctest agent implement` creates a new session.

**Setup**: Mock config with one text event. Invoke `doctest agent implement "implement feature"`.

**RED expectation**: `doctest agent implement` sub-command doesn't exist → exit non-zero.

**GREEN expectation**:
- Exit code 0
- `CODEX_THREAD_ID` environment variable is set after the call
- Session directory exists under `~/.agent-pro/data/doctest-agent/`

### 2. `first-call/no-prompt-error/`

**What it tests**: Calling `doctest agent implement` without a prompt fails.

**Setup**: Invoke `doctest agent implement` with no arguments.

**RED expectation**: Sub-command doesn't exist → exit non-zero (same as above).

**GREEN expectation**:
- Exit code non-zero
- Stderr contains "requires <prompt>"

### 3. `completion/sub-agent-completes/`

**What it tests**: Sub-agent finishes its work successfully.

**Setup**: Mock config with text event `"I have implemented the feature."`. Invoke with a prompt.

**RED expectation**: Sub-command doesn't exist → exit non-zero.

**GREEN expectation**:
- Exit code 0
- Stdout contains `"I have implemented the feature."`

### 4. `questions/sub-agent-yields-questions/`

**What it tests**: Sub-agent calls `yield-pending-questions` and the command exits with questions.

**Setup**: Mock config with a `before_exit` hook that runs `yield-pending-questions '{"id":"1","question":"What port?"}'`. Invoke with a prompt.

**RED expectation**: `yield-pending-questions` binary doesn't exist in PATH → hook fails, or sub-command doesn't exist.

**GREEN expectation**:
- Exit code 0 (questions are not an error — main agent re-invokes after answering)
- Stdout contains the question `"What port?"`

### 5. `questions/resume-after-questions/`

**What it tests**: Setting `CODEX_THREAD_ID` before calling `doctest agent implement` routes to the existing session.

**Setup**: Pre-set `CODEX_THREAD_ID` in env (simulating a previous call). Mock config with a text event `"resumed and done"`. Invoke with answers.

**RED expectation**: Sub-command doesn't exist → exit non-zero.

**GREEN expectation**:
- Exit code 0
- Stdout contains `"resumed and done"`
- The mock's `--session` flag matches `CODEX_THREAD_ID` value

### 6. `failure/sub-agent-exits-nonzero/`

**What it tests**: Sub-agent fails and the error propagates.

**Setup**: Mock config with `exit_code: 3` and `stderr: "implementation failed"`. Invoke with a prompt.

**RED expectation**: Sub-command doesn't exist → exit non-zero (but for wrong reason).

**GREEN expectation**:
- Exit code non-zero
- Stderr or stdout contains `"implementation failed"` or the error message

### 7. `yield-pending-questions-binary/writes-to-fifo/`

**What it tests**: The `yield-pending-questions` binary (dispatched from the doctest binary) writes questions to the configured FIFO.

**Setup**: Create a named FIFO pipe. Set `QUESTION_FIFO` env var. Invoke the binary directly (not via `agent implement`) as `./doctest-bin` renamed to `yield-pending-questions` with question JSON args. Read from the FIFO.

**RED expectation**: Binary doesn't dispatch `yield-pending-questions` → error.

**GREEN expectation**:
- FIFO contains valid JSON line with `"type":"question"`, `"id":"1"`, `"question":"What port?"`

### 8. `yield-pending-questions-binary/no-fifo-error/`

**What it tests**: Calling `yield-pending-questions` without `QUESTION_FIFO` set fails cleanly.

**Setup**: Invoke the binary (as `yield-pending-questions`) without `QUESTION_FIFO` in the environment.

**RED expectation**: Binary doesn't dispatch → error.

**GREEN expectation**:
- Exit code non-zero
- Stderr contains `"QUESTION_FIFO must be set"`

---

## RED Phase: What Fails and Why

At RED phase, `doctest agent implement` doesn't exist at all. Expected failures:

| Test | RED Failure |
|------|------------|
| `creates-session` | `unknown agent command: implement` |
| `no-prompt-error` | `unknown agent command: implement` |
| `sub-agent-completes` | `unknown agent command: implement` |
| `sub-agent-yields-questions` | `unknown agent command: implement` |
| `resume-after-questions` | `unknown agent command: implement` |
| `sub-agent-exits-nonzero` | `unknown agent command: implement` |
| `writes-to-fifo` | Binary dispatch not implemented |
| `no-fifo-error` | Binary dispatch not implemented |

All tests should fail in a consistent, predictable way — confirming RED before implementation.

---

## Sealing

After confirming RED:

```sh
git add agents/doctest/libdoc/agent/tests/
```

---

## GREEN Phase: Implementation Order

1. Add `yield-pending-questions` dispatch to `cli.go` (satisfies binary tests)
2. Add `agent implement` sub-command shell to `cli.go` (satisfies "unknown command" tests)
3. Implement session management in `implement.go` (satisfies session tests)
4. Implement agent runner call (satisfies completion/failure tests)
5. Implement question protocol (satisfies yield/resume tests)
