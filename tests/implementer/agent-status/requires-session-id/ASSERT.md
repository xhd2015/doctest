## Expected
- Stderr contains a message about `--session-id` being required.
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

    if !strings.Contains(resp.Stderr, "session") {
        t.Fatalf("stderr missing 'session', got:\n%s", resp.Stderr)
    }
}
```
