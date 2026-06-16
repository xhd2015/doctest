# Scenario

**Feature**: a doctest binary is available at `req.Bin` (built by the root SETUP.md)

```
# logf formats agent session events for display
doctest agent logf <session-id> -> reads event file -> formatted text -> stdout

# show-status reports session progress
doctest agent show-status <session-id> -> session state -> stdout
```

## Preconditions
- A doctest binary is available at `req.Bin` (built by the root SETUP.md).
- Session directories are created under a temp `DOCTEST_DEBUG_SESSION_HOME`.

## Steps
1. Each leaf creates a session directory with `meta.json` and optionally `events.jsonl`.
2. The root `Run` shells out to `doctest agent implement --status --session-id X`.
3. Leaves assert the output format.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=show-status")
    return nil
}
```
