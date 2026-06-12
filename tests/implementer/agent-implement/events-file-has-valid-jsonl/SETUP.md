## Preconditions
- A `CODEX_THREAD_ID` is provided for a new session.
- The mock config contains multiple events (started, updated, completed).

## Steps
1. Set `CODEX_THREAD_ID`.
2. Write mock config with multiple events.
3. Run `doctest agent implement "implement the feature"`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_events_valid")
    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","stdout_events":[{"type":"item.started","item":{"id":"r1","type":"reasoning"}},{"type":"item.completed","item":{"id":"r1","type":"reasoning","text":"thinking...","status":"completed"}},{"type":"item.completed","item":{"id":"m1","type":"message","text":"done","status":"completed"}}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement the feature"}
    return nil
}
```
