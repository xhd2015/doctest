---
label: heavy
---

## Expected
- Both runs exit 0.
- First captured run is cache-hit.
- Second run remains cache-hit: prose-only SETUP.md change does not change spine Go
  (leaf-cache key extracts Go blocks only).

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
    if req.MRFirst == nil || req.MRSecond == nil {
        t.Fatal("multi-run state not set on req (doMultiRun)")
    }
    if !stdoutHasPositiveCached(req.MRFirst.Stdout) {
        t.Fatalf("first run was not cached; stdout:\n%s", req.MRFirst.Stdout)
    }
    if !stdoutHasPositiveCached(req.MRSecond.Stdout) {
        t.Fatalf("second run lost cache after prose-only SETUP edit; expected Cached > 0:\n%s", req.MRSecond.Stdout)
    }
}
```
