# Scenario

**Feature**: tests in this group run the `yield-pending-questions` binary directly

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- Tests in this group run the `yield-pending-questions` binary directly.
- The binary is dispatched via the doctest binary copied as `yield-pending-questions`.

## Steps
1. Read `YIELD_PQ_BIN` from environment.
2. Set the binary for execution.

```go
import (
    "os"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=yield-pending-questions")
    req.Bin = req.YieldPQBin
    return nil
}
```
