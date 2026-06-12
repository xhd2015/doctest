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
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"done","status":"completed"}}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--session-id", "my-sess-flag", "--agent-runner", "fake-codex", "test"}
    return nil
}
```
