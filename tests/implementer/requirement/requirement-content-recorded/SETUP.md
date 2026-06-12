## Preconditions
- A requirement file exists with specific content.

## Steps
1. Set CODEX_THREAD_ID for deterministic session lookup.
2. Write a requirement file containing "feature spec here".
3. Write a mock config for fake-codex.
4. Run `doctest agent implement --agent-runner fake-codex --requirement <file>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env,
		"CODEX_THREAD_ID=codex-tid-334",
	)

	reqFile := writeRequirementFile(t, req, "feature spec here")

	writeMockConfig(t, req, `{
		"version":"agent-pro.fake-runner.v1",
		"runner":"fake-codex",
		"session_id":"inner-session-334",
		"stdout_events":[
			{"type":"item.completed","item":{"id":"m1","type":"message","text":"done","status":"completed"}}
		]
	}`)

	req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "--requirement", reqFile}
	return nil
}
```
