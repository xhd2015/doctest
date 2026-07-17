---
label: heavy
---

## Expected
- Exit code 0.
- Trace numbering is continuous: [1], [2], [3] for the 3 displayable events.
- No gaps in numbering (e.g., [4] should not appear since there are only 3 displayable events).
- Each number appears next to the correct event content.
- System events (step_start) do NOT appear in the trace output.

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

    if !strings.Contains(stdout, "[1]") {
        t.Errorf("stdout missing trace number [1], got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "[2]") {
        t.Errorf("stdout missing trace number [2], got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "[3]") {
        t.Errorf("stdout missing trace number [3], got:\n%s", stdout)
    }

    if strings.Contains(stdout, "[4]") {
        t.Errorf("stdout contains [4] but only 3 displayable events — numbering has gaps, got:\n%s", stdout)
    }

    if strings.Contains(stdout, "[5]") {
        t.Errorf("stdout contains [5] but only 3 displayable events — non-displayable events consumed numbers, got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "First displayable message") {
        t.Errorf("stdout missing event content 'First displayable message', got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "Build project") {
        t.Errorf("stdout missing event content 'Build project', got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "Final output message") {
        t.Errorf("stdout missing event content 'Final output message', got:\n%s", stdout)
    }

    if t.Failed() {
        t.Fatalf("full stdout:\n%s\nstderr:\n%s", stdout, resp.Stderr)
    }
}
```
