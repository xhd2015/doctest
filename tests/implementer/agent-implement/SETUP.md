# Scenario

**Feature**: tests in this group test `doctest agent implement`

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- Tests in this group test `doctest agent implement`.
- The doctest binary and fake-codex are already built by parent.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=agent-implement")
    return nil
}
```
