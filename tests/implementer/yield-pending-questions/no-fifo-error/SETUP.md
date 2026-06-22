# Scenario

**Feature**: `QUESTION_FIFO` environment variable is NOT set

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- `QUESTION_FIFO` environment variable is NOT set.

## Steps
1. Invoke yield-pending-questions without `QUESTION_FIFO`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{`{"id":"1","question":"test"}`}
    req.Env = append(req.Env, "QUESTION_FIFO=")
    return nil
}
```
