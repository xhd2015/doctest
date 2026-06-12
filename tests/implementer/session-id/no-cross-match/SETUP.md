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
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"first call A done","status":"completed"}}
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
