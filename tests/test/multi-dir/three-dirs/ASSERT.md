---
label: heavy
---

## Expected
- Exit code 0.
- Stderr contains test results for `test_a`, `test_b`, and `test_c`.

## Exit Code
- Exit code 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "test_a") {
        t.Fatalf("stderr missing test_a:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "test_b") {
        t.Fatalf("stderr missing test_b:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "test_c") {
        t.Fatalf("stderr missing test_c:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Fatalf("stdout missing PASS:\n%s", resp.Stdout)
    }
}
```
