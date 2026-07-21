---
label: heavy
---

## Expected
- Both runs exit 0.
- First captured run is cache-hit (both leaves warm).
- Second run is `0 Cached` after shared intermediate SETUP.md change
  (spine ancestor of both leaves → both leaf keys miss).

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
        t.Fatal("multi-run state not set")
    }
    if !stdoutHasPositiveCached(state.FirstResp.Stdout) {
        t.Fatalf("first run was not cached; stdout:\n%s", state.FirstResp.Stdout)
    }
    if !strings.Contains(state.SecondResp.Stdout, ", 0 Cached") {
        t.Fatalf("second run was cached after shared mid SETUP.md edit; expected ', 0 Cached':\n%s", state.SecondResp.Stdout)
    }
}
```
