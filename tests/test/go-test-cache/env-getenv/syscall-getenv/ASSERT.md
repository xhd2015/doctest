## Expected
- First captured run reports `1 Cached`.
- Second captured run (different `DOCTEST_SESSION_ID`) still reports `1 Cached`
  because `syscall.Getenv` bypasses `testlog` and is not part of the cache key.

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
    if envState.FirstResp == nil || envState.SecondResp == nil {
        t.Fatal("env-cache state not set by Setup")
    }
    firstStdout := envState.FirstResp.Stdout
    if !strings.Contains(firstStdout, ", 1 Cached") {
        t.Fatalf("first run was not cached; expected stdout to contain ', 1 Cached':\n%s", firstStdout)
    }
    secondStdout := envState.SecondResp.Stdout
    if !strings.Contains(secondStdout, ", 1 Cached") {
        t.Fatalf("second run was not cached with syscall.Getenv; expected ', 1 Cached':\n%s", secondStdout)
    }
}
```