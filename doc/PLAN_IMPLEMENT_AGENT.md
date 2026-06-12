# PLAN: `doctest agent implement`

## Overview

Add a `doctest agent implement "<prompt>"` sub-command that spawns a sub-agent
via the agent runner (opencode by default), blocks until the sub-agent completes
or yields questions, and supports session continuity via `CODEX_THREAD_ID`.

```
doctest agent implement "<initial prompt>"     # first call — creates session
doctest agent implement "<followup/answers>"    # subsequent calls — same session
```

---

## 1. File Structure

```
agents/doctest/
├── libdoc/
│   ├── cli/
│   │   └── cli.go                      # [MODIFY] add "implement" case under agent
│   └── agent/
│       ├── implement.go                # [NEW] core logic for agent implement
│       └── implement_test.go           # [NEW] doctests
```

---

## 2. CLI Changes (`libdoc/cli/cli.go`)

### 2.1 Usage update

Add to the usage string:

```
  agent implement <prompt>
```

### 2.2 Switch case

Add under `case "agent":` → `runAgent()`:

```go
case "implement":
    return runAgentImplement(args[1:])
```

### 2.3 `runAgentImplement`

```go
func runAgentImplement(args []string) error {
    var agentRunnerFlag *string
    remainArgs, err := lessflags.String("--agent-runner", &agentRunnerFlag).
        Help("-h,--help", agentImplementUsage).
        Parse(args)
    if err != nil {
        return err
    }
    prompt := strings.Join(remainArgs, " ")
    agentRunner := "opencode"
    if agentRunnerFlag != nil {
        agentRunner = *agentRunnerFlag
    }
    return agent.Implement(agent.ImplementOptions{
        Prompt:      prompt,
        AgentRunner: agentRunner,
    })
}
```

### 2.4 Help text

```
Usage: doctest agent implement [--agent-runner RUNNER] <prompt>

Spawn a sub-agent to implement code that makes doctests pass.
Blocks until the sub-agent completes or yields questions via
yield-pending-questions.

Options:
  --agent-runner RUNNER   opencode, codex, or fake-codex (default: opencode)
  -h, --help              Show help
```

---

## 3. Core Logic (`libdoc/agent/implement.go`)

### 3.1 Options

```go
type ImplementOptions struct {
    Prompt      string
    AgentRunner string
}
```

### 3.2 Session Management

Use a session directory under `~/.agent-pro/data/doctest-agent/`:

```
~/.agent-pro/data/doctest-agent/
└── <thread-id>/
    ├── meta.json          # { thread_id, created_at, prompt, codepath }
    └── events.jsonl       # all events in the session
```

**Thread ID generation**: `impl_` prefix + timestamp-based UUID.

**CODEX_THREAD_ID env var**: Set before spawning the sub-agent. On first call,
generate a new ID. On subsequent calls, read from env.

### 3.3 `func Implement(opts ImplementOptions) error`

```go
func Implement(opts ImplementOptions) error {
    // 1. Detect or create session
    threadID := os.Getenv("CODEX_THREAD_ID")
    sessionDir := resolveSessionDir(threadID)

    if threadID == "" {
        // First call — create new session
        threadID = newThreadID()
        sessionDir = createSessionDir(threadID)
        os.Setenv("CODEX_THREAD_ID", threadID)
        writeMeta(sessionDir, opts)
    }

    // 2. Prepare runtime environment
    tempDir := os.MkdirTemp("", "doctest-agent-*")
    defer os.RemoveAll(tempDir)

    // Copy self as yield-pending-questions binary
    copyYieldPendingQuestions(tempDir)
    prependToPath(tempDir)

    // 3. Set up question protocol files
    questionFifo := filepath.Join(tempDir, "question.fifo")
    syscall.Mkfifo(questionFifo, 0666)
    answerDir := filepath.Join(tempDir, "answer")
    os.Mkdir(answerDir, 0755)

    os.Setenv("QUESTION_FIFO", questionFifo)
    os.Setenv("ANSWER_DIR", answerDir)

    // 4. Start question monitor goroutine
    questionCh := make(chan Question, 16)
    go monitorQuestions(questionFifo, questionCh)

    // 5. Build the prompt
    //    On first call: the raw prompt with context about the test tree
    //    On subsequent calls: format as followup to existing session
    fullPrompt := buildPrompt(opts, threadID)

    // 6. Call the agent runner (blocking)
    output, err := runAgent(fullPrompt, opts.AgentRunner, threadID, sessionDir)

    // 7. Check if questions were yielded
    select {
    case q := <-questionCh:
        // Sub-agent yielded questions — print them and return
        printQuestions(q)
        return nil  // exit non-error so main agent can process
    default:
    }

    // 8. No questions — sub-agent completed
    if err != nil {
        return fmt.Errorf("sub-agent failed: %w", err)
    }
    fmt.Println(output)
    return nil
}
```

### 3.4 Question Protocol

Follows the `add-pending-questions` pattern:

**Sub-agent side** (`yield-pending-questions` binary):
- Reads `QUESTION_FIFO` env var
- Accepts questions as JSON args
- Writes each question as a JSON line to the named FIFO
- Exits immediately (does not wait for answers)

**Main agent side** (monitor goroutine):
- Opens the FIFO for reading
- Reads JSON lines, each a question
- Sends questions to a channel
- Main `Implement` func detects questions after `runAgent` returns

**Question format**:
```json
{"type":"question","id":"1","question":"What error format is expected?","options":[{"label":"JSON","description":"..."}]}
```

### 3.5 Agent Runner Call

```go
func runAgent(prompt, agentRunner, threadID, sessionDir string) (string, error) {
    env := agentexec.NewEnv(&agentexec.PathsConfig{
        RootDirName: ".agent-pro",
        DataDirName: "data",
        BinDirName:  "bin",
    }, "AGENT_PRO_CONFIG_HOME")

    runner, err := agentprovider.Build(registry.AgentRunnerID(agentRunner), "", ".", env)
    if err != nil {
        return "", err
    }

    opts := &registry.AskOptions{
        Model:     "",  // use default
        Workspace: ".",
        SessionID: threadID,
    }

    return runner.Agent.Ask(context.Background(), prompt, opts, nil)
}
```

### 3.6 Prompt Construction

**First call**:
```
## Task
Implement the feature to make all doctests pass.

## Test Tree
<list of test leaves and what they cover>

## Constraints
- Tests are sealed via git add — do not modify test files
- Call `yield-pending-questions` if you need clarification

## Implementation Request
<prompt>
```

**Subsequent calls** (followup):
```
<answers or followup message>
```
Sent to the existing session; the agent runner handles thread continuity.

---

## 4. `yield-pending-questions` Binary

Follows the `add-pending-questions` unified binary dispatch pattern.

### 4.1 Dispatch

In `Implement()`, copy the current binary to `tempDir/yield-pending-questions`.
Prepend `tempDir` to `PATH`. The agent runner's sub-agent can then call
`yield-pending-questions` directly.

### 4.2 Implementation

```go
// In implement.go or a separate file
func handleYieldPendingQuestions() {
    questionFifo := os.Getenv("QUESTION_FIFO")
    args := os.Args[1:]

    // Each arg is a JSON question object
    fifo, err := os.OpenFile(questionFifo, os.O_WRONLY, 0)
    for _, arg := range args {
        // Parse JSON, wrap in question entry, write to FIFO
        entry := QuestionEntry{Type: "question", ...}
        json.NewEncoder(fifo).Encode(entry)
    }
    fifo.Close()
}
```

### 4.3 CLI dispatch

In `main.go` (or `cli.go`), before the normal dispatch:

```go
if filepath.Base(os.Args[0]) == "yield-pending-questions" {
    handleYieldPendingQuestions()
    return
}
```

---

## 5. `cli.go` — Add dispatch hook

At the top of `Run()` in `cli.go`:

```go
func Run(args []string) error {
    // Dispatch: if invoked as yield-pending-questions, run that logic
    if filepath.Base(os.Args[0]) == "yield-pending-questions" {
        return agent.HandleYieldPendingQuestions(args)
    }
    // ... existing dispatch
}
```

---

## 6. Test Plan (Doctests)

Tests live in `libdoc/agent/implement_test.go` and use `fake-codex` /
`fake-opencode` via `--mock-config` — no real LLM calls.

### Test Setup Pattern

Each test:
1. Creates a valid doctest tree in a temp dir (with RED-failing stubs)
2. Runs `git init && git add <test-dir>` to simulate sealed tests
3. Creates a mock config JSON for fake-opencode
4. Sets `FAKE_OPENCODE_MOCK_CONFIG` env var (or `FAKE_CODEX_MOCK_CONFIG`)
5. Calls `agent.Implement(opts)` or invokes the CLI
6. Asserts expected output / behavior

### Test Cases

#### Test 1: First call creates session
- Mock returns a simple completion event
- Asserts: command exits 0, `CODEX_THREAD_ID` is set, session dir exists

#### Test 2: Sub-agent completes successfully
- Mock returns events simulating code generation
- Asserts: command exits 0, output contains expected text

#### Test 3: Sub-agent yields questions
- Mock uses a hook (`before_exit`) that calls `yield-pending-questions`
- Hook calls: `yield-pending-questions '{"id":"1","question":"What is the port?"}'`
- Asserts: command exits 0, question appears in output

#### Test 4: Re-invoke after questions (session continuity)
- First call yields questions (as in Test 3)
- Second call with same `CODEX_THREAD_ID` provides answers
- Mock for second call completes successfully
- Asserts: second call exits 0, session is same thread

#### Test 5: Sub-agent fails
- Mock returns non-zero exit code
- Asserts: command returns error

#### Test 6: verify path passing to sub-agent
- Verify that working directory and test tree are accessible to sub-agent
- Mock's hook reads the test tree and confirms files exist

### Mock Config Example

```json
{
  "runner": "fake-opencode",
  "session_id": "fake-session",
  "stdout_events": [
    {
      "type": "text",
      "part": {
        "id": "p1",
        "type": "text",
        "text": "I have implemented the feature."
      }
    }
  ],
  "hooks": [
    {
      "at": "before_exit",
      "event": "yield_questions",
      "payload": {
        "questions": [
          {"id": "1", "question": "What port?"}
        ]
      }
    }
  ],
  "hook_command": "yield-pending-questions '{\"id\":\"1\",\"question\":\"What port?\"}'"
}
```

### Test Execution

```sh
cd agents/doctest
go test ./libdoc/agent/ -run TestImplement -v
```

---

## 7. Implementation Order

1. **`implement.go`** — core `Implement()` function + session management
2. **`implement.go`** — `yield-pending-questions` binary handler
3. **`cli.go`** — add `implement` sub-command + dispatch hook
4. **`implement_test.go`** — doctests (using fake-opencode)
5. **`doc.go`** — embed the spec doc (optional, for `skill show`)

---

## 8. Edge Cases

| Scenario | Behavior |
|----------|----------|
| No prompt provided | Error: `agent implement requires <prompt>` |
| `CODEX_THREAD_ID` set but session dir missing | Treat as first call (create new session) |
| `yield-pending-questions` called with no FIFO | Error: `QUESTION_FIFO must be set` |
| Sub-agent modifies staged test files | Main agent detects via `git diff` (Phase 7) |
| Multiple question rounds | Each round: sub-agent yields → main agent answers → re-invoke → loop |
| CTRL+C during blocked `Implement` | Clean up temp dir, preserve session for resume |
