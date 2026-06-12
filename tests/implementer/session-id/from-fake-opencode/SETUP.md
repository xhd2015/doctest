## Preconditions
- Both `DOCTEST_AGENT_IMPLEMENTER_SESSION_ID` and `CODEX_THREAD_ID` are set.
- The env var takes priority over CODEX_THREAD_ID.

## Steps
1. Set `DOCTEST_AGENT_IMPLEMENTER_SESSION_ID=oc-sess-111` (Priority 1).
2. Set `CODEX_THREAD_ID=should-be-ignored` (proves env var wins).
3. Write a mock config for fake-codex.
4. Run `doctest agent implement "test" --agent-runner fake-codex`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env,
        "DOCTEST_AGENT_IMPLEMENTER_SESSION_ID=oc-sess-111",
        "CODEX_THREAD_ID=should-be-ignored",
    )

    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-from-env",
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"done","status":"completed"}}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "test"}
    return nil
}
```
