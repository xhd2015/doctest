# Scenario

**Feature**: `--session-id` flag is passed

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

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
        "llm_events":[
            {"type":"message","text":"done"}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--session-id", "my-print-sess", "--agent-runner", "fake-codex", "test"}
    return nil
}
```
