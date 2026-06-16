# Scenario

**Feature**: tests in this group verify that `--trace` numbering is continuous (no gaps caused by non-displayable events)

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- Tests in this group verify that `--trace` numbering is continuous (no gaps caused by non-displayable events).
- The fix moves `n++` inside the `if formatted != ""` block so only visible output consumes a number.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=agent-trace-numbering")
    return nil
}
```
