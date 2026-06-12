## Preconditions
- A `CODEX_THREAD_ID` is provided for a new session.
- The mock config contains known event IDs that can be matched.

## Steps
1. Set `CODEX_THREAD_ID`.
2. Write mock config with known events.
3. Run `doctest agent implement "implement the feature"`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_events_match")
    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","stdout_events":[{"type":"item.started","item":{"id":"evt_a","type":"reasoning"}},{"type":"item.completed","item":{"id":"evt_a","type":"reasoning","text":"thinking...","status":"completed"}},{"type":"item.completed","item":{"id":"evt_b","type":"message","text":"all done","status":"completed"}}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement the feature"}
    return nil
}
```
