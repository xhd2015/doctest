## Preconditions
- No session ID env vars are set.
- No opencode ancestor exists in the process tree (doctest runner is the parent).

## Steps
1. Set no env vars at all.
2. Write a mock config for fake-codex.
3. Run `doctest agent implement "test" --agent-runner fake-codex`.
4. Expect the implementer to fail with a clear error message.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    writeMockConfig(t, req, `{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-missing",
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"done","status":"completed"}}
        ]
    }`)

    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "test"}
    return nil
}
```
