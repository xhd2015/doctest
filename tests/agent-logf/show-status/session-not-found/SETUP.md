## Preconditions
- No session directory exists for the given session ID.

## Steps
1. Set `DOCTEST_DEBUG_SESSION_HOME` to an empty temp directory.
2. Run `doctest agent implement --status --session-id test-status-missing`.
3. Verify stderr contains an error message without timestamp prefix.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    sessionHome := t.TempDir()
    req.Env = append(req.Env, "DOCTEST_DEBUG_SESSION_HOME="+sessionHome)
    req.Args = []string{"agent", "implement", "--status", "--session-id", "test-status-missing"}
    return nil
}
```
