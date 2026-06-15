## Preconditions
- A unique session ID is provided via `--session-id` flag.

## Steps
1. Run `doctest agent implement "test" --session-id sess-display-create-1 --agent-runner fake-codex`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-sess-1",
        "llm_events":[
            {"type":"message","text":"done"}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--session-id", "sess-display-create-1", "--agent-runner", "fake-codex", "test"}
    return nil
}
```
