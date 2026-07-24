# Scenario

**Feature**: no session directory exists for the given session ID

```
# logf formats agent session events for display
doctest agent logf <session-id> -> reads event file -> formatted text -> stdout

# show-status reports session progress
doctest agent show-status <session-id> -> session state -> stdout
```

## Preconditions
- No session directory exists for the given session ID.

## Steps
1. Set `DOCTEST_DEBUG_SESSION_HOME` to an empty temp directory.
2. Run `doctest agent implement --status --session-id test-status-missing`.
3. Verify stderr contains an error message without timestamp prefix.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    sessionHome := t.TempDir()
    req.Env = append(req.Env, "DOCTEST_DEBUG_SESSION_HOME="+sessionHome)
    req.Args = []string{"agent", "implement", "--status", "--session-id", "test-status-missing"}
    return nil
}
```
