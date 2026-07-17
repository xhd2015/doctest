---
label: heavy
---

## Expected
- Exit code 0.
- Stdout contains `No events` or similar.
- Stdout contains session ID `status-no-events-test`.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    stdout := resp.Stdout

    if !strings.Contains(strings.ToLower(stdout), "no events") {
        t.Fatalf("stdout missing 'no events', got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "status-no-events-test") {
        t.Fatalf("stdout missing session id 'status-no-events-test', got:\n%s", stdout)
    }
}
```
