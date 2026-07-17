---
label: heavy
---

## Expected
- Exit code 0.
- Header shows `Events: 8 lines`.
- Event listing shows only last 3 events (6, 7, 8): `Sixth event`, `Seventh event text here`, `Eighth event - final output`.
- Early event content `early-event-1` must NOT appear in stdout.
- Numbering must be absolute (matching trace positions): [6], [7], [8] NOT [1], [2], [3].

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

    if strings.Contains(stdout, "[1]") {
        t.Errorf("stdout contains [1] but should use absolute numbering [6],[7],[8], got:\n%s", stdout)
    }

    if strings.Contains(stdout, "[2]") {
        t.Errorf("stdout contains [2] but should use absolute numbering [6],[7],[8], got:\n%s", stdout)
    }

    if strings.Contains(stdout, "[3]") {
        t.Errorf("stdout contains [3] but should use absolute numbering [6],[7],[8], got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "[6]") {
        t.Errorf("stdout missing absolute number [6] for sixth event, got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "[7]") {
        t.Errorf("stdout missing absolute number [7] for seventh event, got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "[8]") {
        t.Errorf("stdout missing absolute number [8] for eighth event, got:\n%s", stdout)
    }

    if t.Failed() {
        t.Fatalf("full stdout:\n%s\nstderr:\n%s", stdout, resp.Stderr)
    }
}
```
