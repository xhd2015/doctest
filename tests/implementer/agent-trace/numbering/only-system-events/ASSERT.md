---
label: heavy
---

## Expected
- Exit code 0.
- Events header shows event count.
- No trace numbers like [1], [2], [3] appear in output (no displayable events).

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

    if strings.Contains(stdout, "[1]") {
        t.Errorf("stdout contains [1] but no event should produce visible output, got:\n%s", stdout)
    }

    if strings.Contains(stdout, "[2]") {
        t.Errorf("stdout contains [2] but no event should produce visible output, got:\n%s", stdout)
    }

    if t.Failed() {
        t.Fatalf("full stdout:\n%s\nstderr:\n%s", stdout, resp.Stderr)
    }
}
```
