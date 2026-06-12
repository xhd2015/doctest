## Preconditions
- A session directory exists from a prior call.

## Steps
1. First call: run `agent implement` with `DOCTEST_AGENT_IMPLEMENTER_SESSION_ID=resume-test-555` to create a session.
2. Record the session directory path from the first call.
3. Second call: run with the same ID again.
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
        "DOCTEST_AGENT_IMPLEMENTER_SESSION_ID=resume-test-555",
    )

    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-resume-sess",
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"first call done","status":"completed"}}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "first call"}
    req.Env = append(req.Env, "IMPLEMENT_CALL=first")

    // store the bin path so ASSERT can run the second call
    req.Env = append(req.Env, "DOCTEST_BIN_FOR_RESUME="+req.Bin)
    _ = os.Stat
    _ = filepath.Join
    _ = exec.Command
    return nil
}
```
