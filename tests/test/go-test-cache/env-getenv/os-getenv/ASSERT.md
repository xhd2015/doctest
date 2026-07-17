---
label: heavy
---

## Expected
- First captured run (same env as warmup) reports `1 Cached`.
- Second captured run (different `DOCTEST_SESSION_ID`) reports `0 Cached` because
  `os.Getenv` is recorded in the go test cache inputs.

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
    if !strings.Contains(secondStdout, ", 0 Cached") {
        t.Fatalf("second run was cached after os.Getenv env change; expected ', 0 Cached':\n%s", secondStdout)
    }
}
```