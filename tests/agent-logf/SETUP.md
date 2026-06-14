## Preconditions
- The doctest module root is the parent of this test tree (`DOCTEST_ROOT/..`).
- Each leaf tests the logging behavior of `Logf`, `traceSession`, or `showStatus`.
- Tests verify that event/status lines use timestamped `Logf` output while UI framing uses bare `fmt.Fprintf`.

## Steps
1. Each grouping node provides its own `Run` function suited to the function under test.
2. Leaves configure inputs via `req.Args` or `req.Env`.

## Context
- These tests verify the contract: Logf produces `[2006-01-02T15:04:05]` prefixed output.
- `traceSession` and `showStatus` tests create real session directories and run the doctest binary.

```go
import (
    "testing"
    "time"
)

type Request struct {
    Args    []string
    Env     []string
    WorkDir string
    Timeout time.Duration
    Bin     string
}

type Response struct {
    ExitCode int
    Stdout   string
    Stderr   string
    Err      error
}

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_FEATURE=agent-logf")
    req.Timeout = 30 * time.Second
    return nil
}
```
