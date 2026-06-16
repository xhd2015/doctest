# Scenario

**Feature**: `--list-sessions` and `--session-id` are mutually exclusive

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- `--list-sessions` and `--session-id` are mutually exclusive.
- When both are given, stderr gets an error and exit code is 0.

## Steps
1. Run `doctest agent implement --list-sessions --session-id some-id`.

```go
func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "implement", "--list-sessions", "--session-id", "some-id"}
    return nil
}
```
