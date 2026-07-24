---
label: heavy
---

## Expected
- Exit code 0 (absolute `/...` path is expanded, not stat'd literally).
- stderr contains "alpha" (tests under alpha/ found).
- stderr does NOT contain "beta" (tests under beta/ excluded).
- stdout contains PASS.

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
        t.Fatalf("exit code = %d, want 0 for absolute /... path; stderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "no such file or directory") && strings.Contains(resp.Stderr, "/...") {
        t.Fatalf("doctest treated /... literally instead of expanding pattern; stderr:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "alpha") {
        t.Fatalf("stderr missing alpha:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "beta") {
        t.Fatalf("stderr should not contain beta:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Fatalf("stdout missing PASS:\n%s", resp.Stdout)
    }
}
```