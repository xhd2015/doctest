## Preconditions
- The mock config returns a text event indicating success.

## Steps
1. Write mock config with a completion event.
2. Run `doctest agent implement "implement feature"`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_completes")
    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","stdout_events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"I have implemented the feature.","status":"completed"}}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement feature"}
    return nil
}
```
