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
