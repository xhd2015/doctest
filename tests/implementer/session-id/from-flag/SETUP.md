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
- `CODEX_THREAD_ID` is also set (should be stored for traceability but NOT used as session ID).

## Steps
1. Run `doctest agent implement "test" --session-id my-sess-flag --agent-runner fake-codex`.
2. Set `CODEX_THREAD_ID=should-be-ignored`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env,
        "CODEX_THREAD_ID=should-be-ignored",
    )

    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-from-flag",
        "llm_events":[
            {"type":"message","text":"done"}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--session-id", "my-sess-flag", "--agent-runner", "fake-codex", "test"}
    return nil
}
```
