# Scenario

**Feature**: first call uses `--session-id sess-A`, second call uses `CODEX_THREAD_ID=sess-B`

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- First call uses `--session-id sess-A`, second call uses `CODEX_THREAD_ID=sess-B`.
- Different sources should NOT cross-match.

## Steps
1. First call: run with `--session-id sess-A`.
2. Second call: run with `CODEX_THREAD_ID=sess-B` (different source, different ID).
3. Two separate sessions should exist.

```go
import (
    "os"
    "os/exec"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-cross-A",
        "llm_events":[
            {"type":"message","text":"first call A done"}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--session-id", "sess-A", "--agent-runner", "fake-codex", "first call"}
    req.Env = append(req.Env, "DOCTEST_BIN_FOR_RESUME="+req.Bin)
    _ = os.Stat
    _ = filepath.Join
    _ = exec.Command
    return nil
}
```
