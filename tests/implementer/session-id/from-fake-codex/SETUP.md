## Preconditions
- `DOCTEST_AGENT_IMPLEMENTER_SESSION_ID` is NOT set.
- `CODEX_THREAD_ID` is set.

## Steps
1. Set `CODEX_THREAD_ID=codex-tid-222` (Priority 2 fallback).
2. Write a mock config for fake-codex.
3. Run `doctest agent implement "test" --agent-runner fake-codex`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env,
        "CODEX_THREAD_ID=codex-tid-222",
    )

    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-from-codex",
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"done","status":"completed"}}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "test"}
    return nil
}
```
