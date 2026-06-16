# Scenario

**Feature**: multiple session ID sources are set simultaneously

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- Multiple session ID sources are set simultaneously.

## Steps
1. Set `DOCTEST_AGENT_IMPLEMENTER_SESSION_ID=store-prio-1`.
2. Set `CODEX_THREAD_ID=store-codex-2`.
3. Write a mock config for fake-codex.
4. Run `doctest agent implement "test" --agent-runner fake-codex`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env,
        "DOCTEST_AGENT_IMPLEMENTER_SESSION_ID=store-prio-1",
        "CODEX_THREAD_ID=store-codex-2",
    )

    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-store-sess",
        "llm_events":[
            {"type":"message","text":"done"}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "test"}
    return nil
}
```
