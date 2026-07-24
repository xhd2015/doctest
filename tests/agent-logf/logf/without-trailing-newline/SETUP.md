# Scenario

**Feature**: the format string `"hello"` has no trailing newline

```
# logf formats agent session events for display
doctest agent logf <session-id> -> reads event file -> formatted text -> stdout

# show-status reports session progress
doctest agent show-status <session-id> -> session state -> stdout
```

## Preconditions
- The format string `"hello"` has no trailing newline.

## Steps
1. Set `LOGF_FORMAT=hello` via env.
2. Call `Logf("hello")` via the parent Run.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Env = append(req.Env, "LOGF_FORMAT=hello")
    return nil
}
```
