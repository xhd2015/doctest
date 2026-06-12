## Preconditions
- The `--agent-runner` flag is forwarded to the sub-agent.

## Steps
1. Set up a mock config.
2. Run `doctest agent implement "task" --agent-runner fake-codex`.
3. Verify the mock config was used (proving the runner flag worked).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_orch_args")
    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","session_id":"sess_args","stdout_events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"args forwarded ok","status":"completed"}}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "test args forwarded"}
    return nil
}
```
