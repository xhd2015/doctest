## Preconditions
- A `CODEX_THREAD_ID` is provided for a new session.
- The mock config contains a simple completion event.

## Steps
1. Set `CODEX_THREAD_ID`.
2. Write mock config.
3. Run `doctest agent implement "implement the feature"`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_new_session")
    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","stdout_events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"implementation done","status":"completed"}}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement the feature"}
    return nil
}
```
