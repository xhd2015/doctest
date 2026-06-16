# Scenario

**Feature**: a unique session ID is provided via `--session-id` flag

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

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
