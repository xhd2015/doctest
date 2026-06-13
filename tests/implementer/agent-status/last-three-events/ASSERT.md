## Expected
- Exit code 0.
- Header shows `Events: 8 lines`.
- Event listing shows only last 3 events (6, 7, 8): `Sixth event`, `Seventh event text here`, `Eighth event - final output`.
- Early event content `early-event-1` must NOT appear in stdout.

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

    if !strings.Contains(stdout, "Events:  8 lines") {
        t.Errorf("header missing event count 8, got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "Sixth event") {
        t.Errorf("stdout missing sixth event, got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "Seventh event text here") {
        t.Errorf("stdout missing seventh event, got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "Eighth event - final output") {
        t.Errorf("stdout missing eighth event, got:\n%s", stdout)
    }

    if strings.Contains(stdout, "early-event-1") {
        t.Errorf("stdout contains truncated event (early-event-1 should not appear), got:\n%s", stdout)
    }

    if t.Failed() {
        t.Fatalf("full stdout:\n%s\nstderr:\n%s", stdout, resp.Stderr)
    }
}
```
