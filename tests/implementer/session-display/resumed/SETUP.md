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
- The first call creates the session; the second call should resume it.

## Steps
1. First call: run `doctest agent implement "first" --session-id sess-display-resume-1 --agent-runner fake-codex`.
2. Record the doctest binary path for the second call in ASSERT.md.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-sess-resume-1",
        "llm_events":[
            {"type":"message","text":"first call done"}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--session-id", "sess-display-resume-1", "--agent-runner", "fake-codex", "first"}
    req.Env = append(req.Env, "DOCTEST_BIN_FOR_RESUME="+req.Bin)
    return nil
}
```
