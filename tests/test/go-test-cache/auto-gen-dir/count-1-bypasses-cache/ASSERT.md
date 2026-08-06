## Expected
- Both runs exit 0.
- First measured run (default flags) is leaf-cached (`Cached` > 0).
- Second measured run with `-count=1` reports `, 0 Cached` (leaf-cache bypass).

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
    firstStdout := req.MRFirst.Stdout
    if !stdoutHasPositiveCached(firstStdout) {
        t.Fatalf("first run (default) was not cached; expected Cached > 0:\n%s", firstStdout)
    }
    secondStdout := req.MRSecond.Stdout
    if !strings.Contains(secondStdout, ", 0 Cached") {
        t.Fatalf("second run was cached with -count=1; expected ', 0 Cached':\n%s", secondStdout)
    }
}
```
