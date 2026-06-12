## Preconditions
- `CODEX_THREAD_ID` is set.

## Steps
1. First call: run `agent implement` with `CODEX_THREAD_ID=resume-codex-999`.
2. Record the session directory path from the first call.
3. Second call: run with the same `CODEX_THREAD_ID` again.
4. Verify the session directory is reused, not recreated.

```go
import (
    "os"
    "os/exec"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env,
        "CODEX_THREAD_ID=resume-codex-999",
    )

    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-codex-resume",
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"first codex call done","status":"completed"}}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "first call"}
    req.Env = append(req.Env, "DOCTEST_BIN_FOR_RESUME="+req.Bin)
    _ = os.Stat
    _ = filepath.Join
    _ = exec.Command
    return nil
}
```
