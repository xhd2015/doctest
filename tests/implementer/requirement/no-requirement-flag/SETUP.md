## Preconditions
- No --requirement flag is used. Only a CLI prompt is provided.

## Steps
1. Set CODEX_THREAD_ID for deterministic session lookup.
2. Write a mock config for fake-codex.
3. Run `doctest agent implement --agent-runner fake-codex "hello"`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env,
		"CODEX_THREAD_ID=codex-tid-337",
	)

	writeMockConfig(t, req, `{
		"version":"agent-pro.fake-runner.v1",
		"runner":"fake-codex",
		"session_id":"inner-session-337",
		"stdout_events":[
			{"type":"item.completed","item":{"id":"m1","type":"message","text":"done","status":"completed"}}
		]
	}`)

	req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "hello"}
	return nil
}
```
