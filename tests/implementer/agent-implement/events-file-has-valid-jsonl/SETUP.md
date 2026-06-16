# Scenario

**Feature**: a `CODEX_THREAD_ID` is provided for a new session

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

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
    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","llm_events":[{"type":"think","text":"thinking..."},{"type":"message","text":"done"}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement the feature"}
    return nil
}
```
