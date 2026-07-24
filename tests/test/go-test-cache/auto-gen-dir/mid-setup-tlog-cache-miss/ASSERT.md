---
label: heavy
---

## Expected
- Both runs exit 0.
- First run cache-hit; second run **0 Cached** after mid Setup t.Log tag change.

## Exit Code
- Exit code 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if req.MRFirst == nil || req.MRSecond == nil {
        t.Fatal("multi-run state not set on req (doMultiRun)")
    }
    if !stdoutHasPositiveCached(req.MRFirst.Stdout) {
        t.Fatalf("first run was not cached; stdout:\n%s", req.MRFirst.Stdout)
    }
    if !strings.Contains(req.MRSecond.Stdout, ", 0 Cached") {
        t.Fatalf("second run was cached after mid t.Log SETUP edit; expected ', 0 Cached':\n%s", req.MRSecond.Stdout)
    }
}
```
