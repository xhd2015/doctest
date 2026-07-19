---
label: heavy
---

## Expected
- Both runs exit 0.
- Second run stdout (summary line) contains ", 1 Cached".
- Second run completes faster than first (under 5 seconds total).

## Exit Code
- Exit code 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if state.FirstResp == nil || state.SecondResp == nil {
        t.Fatal("multi-run state not set by Run")
    }
    if state.FirstResp.ExitCode != 0 {
        t.Fatalf("first run exit %d, stderr:\n%s", state.FirstResp.ExitCode, state.FirstResp.Stderr)
    }
    secondStdout := state.SecondResp.Stdout
    // Unified suite: one package → ", 1 Cached" (or "N Cached" if leaf counts).
    if !strings.Contains(secondStdout, "Cached") || strings.Contains(secondStdout, "0 Cached") {
        t.Fatalf("second run not cached; expected Cached > 0 in stdout:\n%s", secondStdout)
    }
}
```
