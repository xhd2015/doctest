---
label: heavy
---

## Expected
- Both runs exit 0.
- First captured run is cache-hit.
- Second run **remains** cache-hit after unread intermediate SETUP WorkDir edit
  (Go testcache: linked binary content ID unchanged after DCE).

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
        t.Fatalf("second run lost cache after unread intermediate SETUP edit; expected Cached > 0:\n%s", state.SecondResp.Stdout)
    }
}
```
