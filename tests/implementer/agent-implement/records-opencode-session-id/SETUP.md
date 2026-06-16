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
- The mock config contains a simple completion event.

## Steps
1. Set `CODEX_THREAD_ID`.
2. Write mock config.
3. Run `doctest agent implement "implement feature"`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_record_sid")
    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","session_id":"mock-sid-123","llm_events":[{"type":"message","text":"implementation done"}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement feature"}
    return nil
}
```
