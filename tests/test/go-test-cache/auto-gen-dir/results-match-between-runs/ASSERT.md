---
label: heavy
---

## Expected
- Both runs have the same exit code.
- Both runs pass (exit 0).
- Both runs show pass in output summary.

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
    if req.MRFirst == nil || req.MRSecond == nil {
        t.Fatal("multi-run state not set on req (doMultiRun)")
    }
    if req.MRFirst.ExitCode != req.MRSecond.ExitCode {
        t.Fatalf("exit codes differ: first=%d second=%d\nfirst:\n%s\nsecond:\n%s",
            req.MRFirst.ExitCode, req.MRSecond.ExitCode,
            req.MRFirst.Stderr, req.MRSecond.Stderr)
    }
    if req.MRFirst.ExitCode != 0 {
        t.Fatalf("expected exit 0, got first=%d second=%d",
            req.MRFirst.ExitCode, req.MRSecond.ExitCode)
    }
    if !strings.Contains(req.MRFirst.Stdout, "0 Fail") {
        t.Fatalf("first run missing pass in stdout:\n%s", req.MRFirst.Stdout)
    }
    if !strings.Contains(req.MRSecond.Stdout, "0 Fail") {
        t.Fatalf("second run missing pass in stdout:\n%s", req.MRSecond.Stdout)
    }
}
```
