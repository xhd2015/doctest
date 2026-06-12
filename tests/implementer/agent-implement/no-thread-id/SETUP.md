## Preconditions
- `CODEX_THREAD_ID` is not set in the environment.

## Steps
1. Run `doctest agent implement` without `CODEX_THREAD_ID`, using `fake-codex`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-no-thread",
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"done","status":"completed"}}
        ]
    }`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "some prompt"}
    return nil
}
```
