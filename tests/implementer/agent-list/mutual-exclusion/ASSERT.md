## Expected
- Stderr contains an error about mutual exclusion (e.g. "cannot use --list-sessions with --session-id").
- Exit code is 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    if !strings.Contains(strings.ToLower(resp.Stderr), "list") || !strings.Contains(strings.ToLower(resp.Stderr), "session") {
        t.Fatalf("stderr missing mutual-exclusion error, got:\n%s", resp.Stderr)
    }
}
```
