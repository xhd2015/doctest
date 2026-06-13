## Expected
- Exit code 0.
- Stdout indicates no sessions found (e.g. "No sessions" or "0 sessions").

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    if !strings.Contains(strings.ToLower(resp.Stdout), "no session") && !strings.Contains(resp.Stdout, "0 session") {
        t.Fatalf("stdout missing 'no sessions' / '0 sessions', got:\n%s", resp.Stdout)
    }
}
```
