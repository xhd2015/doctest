## Expected
- Both runs exit 0.
- Second run stderr does NOT contain "(cached)" because `-count=2` bypasses the cache.

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
    firstStdout := state.FirstResp.Stdout
    if !strings.Contains(firstStdout, "(cached)") {
        t.Fatalf("first run (count=1) was not cached; expected stdout to contain '(cached)':\n%s", firstStdout)
    }
    secondStdout := state.SecondResp.Stdout
    if strings.Contains(secondStdout, "(cached)") {
        t.Fatalf("second run was cached with -count=2; expected no cache:\n%s", secondStdout)
    }
}
```
