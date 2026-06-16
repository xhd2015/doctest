# Scenario

**Feature**: `CODEX_THREAD_ID` is set to `flag-fallback-test`

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- `CODEX_THREAD_ID` is set to `flag-fallback-test`.
- No `--session-id` flag on first call.

## Steps
1. First call: run `agent implement` with `CODEX_THREAD_ID=flag-fallback-test`.
2. Record the session directory path from the first call.
3. Second call: run with `--session-id flag-fallback-test` (same value, different source).
4. Verify the session directory is found and reused (not a new session created).

```go
import (
    "os"
    "os/exec"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env,
        "CODEX_THREAD_ID=flag-fallback-test",
    )

    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-fallback-1",
        "llm_events":[
            {"type":"message","text":"first fallback call done"}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "first fallback call"}
    req.Env = append(req.Env, "DOCTEST_BIN_FOR_RESUME="+req.Bin)
    _ = os.Stat
    _ = filepath.Join
    _ = exec.Command
    return nil
}
```
