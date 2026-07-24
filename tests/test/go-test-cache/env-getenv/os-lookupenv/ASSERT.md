---
label: heavy
---

## Expected
- First captured run (env A) reports Cached > 0.
- Second captured run (different env B) still reports Cached > 0 because
  process env is **not** part of the leaf-cache key (`os.LookupEnv` is not special-cased).

## Exit Code
- Exit code 0.

```go
import (
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if req.MRFirst == nil || req.MRSecond == nil {
        t.Fatal("env-cache state not set by Setup (req.MRFirst/MRSecond)")
    }
    firstStdout := req.MRFirst.Stdout
    if !stdoutHasPositiveCached(firstStdout) {
        t.Fatalf("first run was not cached; expected Cached > 0:\n%s", firstStdout)
    }
    secondStdout := req.MRSecond.Stdout
    if !stdoutHasPositiveCached(secondStdout) {
        t.Fatalf("second run lost cache after os.LookupEnv env change; env must not be in leaf-cache key:\n%s", secondStdout)
    }
}
```
