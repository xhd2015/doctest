## Expected
- Exit code 0.
- stderr contains "alpha" (tests under alpha/ found).
- stderr does NOT contain "beta" (tests under beta/ excluded).
- stderr contains PASS.

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
        t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "alpha") {
        t.Fatalf("stderr missing alpha:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "beta") {
        t.Fatalf("stderr should not contain beta:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "PASS") {
        t.Fatalf("stderr missing PASS:\n%s", resp.Stderr)
    }
}
```
