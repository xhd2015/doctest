# Scenario

**Feature**: the requirement file path does not exist

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

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
		"llm_events":[
			{"type":"message","text":"done"}
		]
	}`)

	req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "--requirement", "/nonexistent/file.md"}
	return nil
}
```
