# Scenario

**Feature**: the implementer SETUP.md provides shared infrastructure

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- The implementer SETUP.md provides shared infrastructure.

## Steps
1. Test meta.json field structure and session persistence across calls.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=meta-fields")
    return nil
}
```
