## Expected
- Both runs exit 0.
- Second run stderr does NOT contain "(cached)" because the SETUP.md changed.

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
    firstStderr := state.FirstResp.Stderr
    if !strings.Contains(firstStderr, "(cached)") {
        t.Fatalf("first run was not cached; expected stderr to contain '(cached)':\n%s", firstStderr)
    }
    secondStderr := state.SecondResp.Stderr
    if strings.Contains(secondStderr, "(cached)") {
        t.Fatalf("second run was cached after SETUP.md edit; expected cache miss:\n%s", secondStderr)
    }
}
```
