# Scenario

**Feature**: `CODEX_THREAD_ID` is not set in the environment

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

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
        "llm_events":[
            {"type":"message","text":"done"}
        ]
    }`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "some prompt"}
    return nil
}
```
