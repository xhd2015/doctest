## Preconditions
- A `CODEX_THREAD_ID` is provided for a new session.
- The mock config contains a simple completion event.

## Steps
1. Set `CODEX_THREAD_ID`.
2. Write mock config.
3. Run `doctest agent implement "implement feature"`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_record_sid")
    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","session_id":"mock-sid-123","stdout_events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"implementation done","status":"completed"}}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement feature"}
    return nil
}
```
