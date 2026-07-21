---
label: heavy
---

## Expected
- Both runs exit 0.
- First captured run is cache-hit (`Cached` > 0).
- After sibling-branch intermediate SETUP edit, second run still has `Cached` > 0
  because the peer leaf's spine is unchanged (leaf-cache key is spine-only).

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
        t.Fatalf("second run lost all cache after sibling branch SETUP edit; peer leaf should stay leaf-cached (Cached > 0):\n%s", state.SecondResp.Stdout)
    }
}
```
