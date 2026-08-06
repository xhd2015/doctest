## Expected
- Both runs exit 0.
- Second run stdout contains ", 0 Cached" because the leaf SETUP Go changed
  (leaf-cache spine key differs).

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
    if !strings.Contains(firstStdout, ", 1 Cached") {
        t.Fatalf("first run was not cached; expected stdout to contain ', 1 Cached':\n%s", firstStdout)
    }
    secondStdout := req.MRSecond.Stdout
    if !strings.Contains(secondStdout, ", 0 Cached") {
        t.Fatalf("second run was cached after SETUP.md edit; expected ', 0 Cached':\n%s", secondStdout)
    }
}
```
