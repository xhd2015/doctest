# Scenario

**Feature**: a requirement file exists

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- A requirement file exists.

## Steps
1. Set CODEX_THREAD_ID for deterministic session lookup.
2. Write a requirement file.
3. Write a mock config for fake-codex.
4. Run `doctest agent implement --agent-runner fake-codex --requirement <file>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env,
		"CODEX_THREAD_ID=codex-tid-338",
	)

	reqFile := writeRequirementFile(t, req, "test requirement content")

	writeMockConfig(t, req, `{
		"version":"agent-pro.fake-runner.v1",
		"runner":"fake-codex",
		"session_id":"inner-session-338",
		"llm_events":[
			{"type":"message","text":"done"}
		]
	}`)

	req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "--requirement", reqFile}
	return nil
}
```
