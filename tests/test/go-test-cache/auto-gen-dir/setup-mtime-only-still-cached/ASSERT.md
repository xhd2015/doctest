---
label: heavy
---

## Expected
- Both runs exit 0.
- First and second captured runs are cache-hits after mtime-only touch
  (content-stable generation does not rewrite gen packages).

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
        t.Fatal("multi-run state not set")
    }
    if !stdoutHasPositiveCached(state.FirstResp.Stdout) {
        t.Fatalf("first run was not cached; stdout:\n%s", state.FirstResp.Stdout)
    }
    if !stdoutHasPositiveCached(state.SecondResp.Stdout) {
        t.Fatalf("second run lost cache after mtime-only touch; expected Cached > 0:\n%s", state.SecondResp.Stdout)
    }
}
```
