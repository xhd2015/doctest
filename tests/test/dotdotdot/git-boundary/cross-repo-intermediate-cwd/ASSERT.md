## Expected
- Exit code 0.
- Stderr contains `inner_test` (test tree found via subdir discovery).
- Stderr contains PASS.

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
    if !strings.Contains(resp.Stderr, "inner_test") {
        t.Fatalf("stderr missing inner_test (should be discovered via subdir fallback):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "PASS") {
        t.Fatalf("stderr missing PASS:\n%s", resp.Stderr)
    }
}
```
