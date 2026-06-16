# Scenario

**Feature**: the implementer SETUP.md provides shared infrastructure (binaries and helpers)

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- The implementer SETUP.md provides shared infrastructure (binaries and helpers).

## Steps
1. Test session display messages for created and resumed sessions.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=session-display")
    return nil
}
```
