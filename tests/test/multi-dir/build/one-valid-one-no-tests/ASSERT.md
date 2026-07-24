---
label: heavy
---

## Expected
- Exit code 0 (the `no_tests` dir is silently skipped, `test_a` builds successfully).

## Exit Code
- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
    }
}
```
