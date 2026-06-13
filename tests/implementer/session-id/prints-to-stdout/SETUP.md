## Preconditions
- `--session-id` flag is passed.
- The mock config returns a simple completion event.

## Steps
1. Run `doctest agent implement "test" --session-id my-print-sess --agent-runner fake-codex`.
2. Check that stdout contains "Session ID:" with the session ID and source info.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-print",
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"done","status":"completed"}}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--session-id", "my-print-sess", "--agent-runner", "fake-codex", "test"}
    return nil
}
```
