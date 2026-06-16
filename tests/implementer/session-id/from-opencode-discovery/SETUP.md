# Scenario

**Feature**: neither `--session-id`, `DOCTEST_AGENT_IMPLEMENTER_SESSION_ID`, nor `CODEX_THREAD_ID` is set

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- Neither `--session-id`, `DOCTEST_AGENT_IMPLEMENTER_SESSION_ID`, nor `CODEX_THREAD_ID` is set.
- Auto-discovery has been removed; no session ID can be resolved.

## Steps
1. Set no session ID sources at all.
2. Write a mock config for fake-codex.
3. Run `doctest agent implement "test" --agent-runner fake-codex`.
4. Expect error with a generated session ID and --session-id usage instructions.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-from-discovery",
        "llm_events":[
            {"type":"message","text":"done"}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "test"}
    return nil
}
```
