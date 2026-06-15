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

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if state.FirstResp == nil || state.SecondResp == nil {
        t.Fatal("multi-run state not set")
    }
    if state.FirstResp.ExitCode != state.SecondResp.ExitCode {
        t.Fatalf("exit codes differ: first=%d second=%d\nfirst:\n%s\nsecond:\n%s",
            state.FirstResp.ExitCode, state.SecondResp.ExitCode,
            state.FirstResp.Stderr, state.SecondResp.Stderr)
    }
    if state.FirstResp.ExitCode != 0 {
        t.Fatalf("expected exit 0, got first=%d second=%d",
            state.FirstResp.ExitCode, state.SecondResp.ExitCode)
    }
    if !strings.Contains(state.FirstResp.Stdout, "0 Fail") {
        t.Fatalf("first run missing pass in stdout:\n%s", state.FirstResp.Stdout)
    }
    if !strings.Contains(state.SecondResp.Stdout, "0 Fail") {
        t.Fatalf("second run missing pass in stdout:\n%s", state.SecondResp.Stdout)
    }
}
```
