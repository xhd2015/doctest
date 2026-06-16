# Scenario

**Feature**: `--status` without `--session-id` must error with a clear message

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- `--status` without `--session-id` must error with a clear message.
- Exit code must be 0 per the "always exit 0" requirement.

## Steps
1. Run `doctest agent implement --status` without `--session-id`.

```go
func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "implement", "--status"}
    return nil
}
```
