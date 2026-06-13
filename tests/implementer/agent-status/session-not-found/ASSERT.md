## Expected
- Stderr contains a message indicating the session was not found.
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

    if !strings.Contains(strings.ToLower(resp.Stderr), "not found") && !strings.Contains(strings.ToLower(resp.Stderr), "no session") {
        t.Fatalf("stderr missing 'not found' or 'no session', got:\n%s", resp.Stderr)
    }
}
```
