---
label: heavy
---

## Expected
- First captured run (env A) reports `1 Cached` (or Cached > 0).
- Second captured run (different env B) still reports Cached > 0 because
  process env is **not** part of the leaf-cache key (`os.Getenv` is not special-cased).

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
    if envState.FirstResp == nil || envState.SecondResp == nil {
        t.Fatal("env-cache state not set by Setup")
    }
    firstStdout := envState.FirstResp.Stdout
    if !stdoutHasPositiveCached(firstStdout) {
        t.Fatalf("first run was not cached; expected Cached > 0:\n%s", firstStdout)
    }
    secondStdout := envState.SecondResp.Stdout
    if !stdoutHasPositiveCached(secondStdout) {
        t.Fatalf("second run lost cache after os.Getenv env change; env must not be in leaf-cache key:\n%s", secondStdout)
    }
}
```
