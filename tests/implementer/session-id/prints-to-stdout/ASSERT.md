## Expected
- Exit code 0.
- Stdout contains "Session ID:" with the session ID `my-print-sess`.
- Stdout shows the source of the session ID (`--session-id`).

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    if !strings.Contains(resp.Stdout, "Session ID:") {
        t.Fatalf("stdout missing 'Session ID:', got:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stdout, "my-print-sess") {
        t.Fatalf("stdout missing session ID 'my-print-sess', got:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stdout, "--session-id") {
        t.Fatalf("stdout missing source indicator '--session-id', got:\n%s", resp.Stdout)
    }
}
```
