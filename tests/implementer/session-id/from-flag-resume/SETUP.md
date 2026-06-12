## Preconditions
- `--session-id` flag is used for both calls.

## Steps
1. First call: run `agent implement` with `--session-id resume-flag-777` to create a session.
2. Record the session directory path from the first call.
3. Second call: run with the same `--session-id` again.
4. Verify the session directory is reused, not recreated.

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
        "session_id":"inner-flag-resume",
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"first flag call done","status":"completed"}}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--session-id", "resume-flag-777", "--agent-runner", "fake-codex", "first call"}
    req.Env = append(req.Env, "DOCTEST_BIN_FOR_RESUME="+req.Bin)
    _ = os.Stat
    _ = filepath.Join
    _ = exec.Command
    return nil
}
```
