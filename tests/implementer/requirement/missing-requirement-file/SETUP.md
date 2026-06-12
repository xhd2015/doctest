## Preconditions
- The requirement file path does not exist.

## Steps
1. Set CODEX_THREAD_ID for deterministic session lookup.
2. Write a mock config for fake-codex.
3. Run `doctest agent implement --agent-runner fake-codex --requirement /nonexistent/file.md`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env,
		"CODEX_THREAD_ID=codex-tid-336",
	)

	writeMockConfig(t, req, `{
		"version":"agent-pro.fake-runner.v1",
		"runner":"fake-codex",
		"stdout_events":[
			{"type":"item.completed","item":{"id":"m1","type":"message","text":"done","status":"completed"}}
		]
	}`)

	req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "--requirement", "/nonexistent/file.md"}
	return nil
}
```
