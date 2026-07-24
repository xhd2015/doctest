# Scenario

**Feature**: the format string contains a line break, special characters, and already ends with `\n`

```
# logf formats agent session events for display
doctest agent logf <session-id> -> reads event file -> formatted text -> stdout

# show-status reports session progress
doctest agent show-status <session-id> -> session state -> stdout
```

## Preconditions
- The format string contains a line break, special characters, and already ends with `\n`.

## Steps
1. Set `LOGF_FORMAT=line1\nline2 -- special: !@#$%\n` via env.
2. Call `Logf("line1\nline2 -- special: !@#$%\n")` via the parent Run.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Env = append(req.Env, "LOGF_FORMAT=line1\nline2 -- special: !@#$%\n")
    return nil
}
```
