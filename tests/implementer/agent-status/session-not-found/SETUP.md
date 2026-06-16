# Scenario

**Feature**: a session ID that does not correspond to any session directory

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- A session ID that does not correspond to any session directory.
- Exit code must be 0 per the "always exit 0" requirement.

## Steps
1. Run `doctest agent implement --status --session-id nonexistent-status-test`.

```go
func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "implement", "--status", "--session-id", "nonexistent-status-test"}
    return nil
}
```
