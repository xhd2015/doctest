---
label: heavy
---

## Expected
- Both runs exit 0.
- Second run stdout (summary line) has `Cached` > 0 via leaf-cache skip
  (or whole-package go `(cached)` expanded to N Cached for all N leaves).

## Exit Code
- Exit code 0.

```go
import (
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
    // Leaf-cache product: warm second run must show Cached > 0.
    if !stdoutHasPositiveCached(secondStdout) {
        t.Fatalf("second run not leaf-cached; expected Cached > 0 in stdout:\n%s", secondStdout)
    }
}
```
