## Expected
- Exit code 0.
- Stderr contains test results for both `alpha_test` and `beta_test`.
- Stderr does NOT contain `hidden_test` (nested go.mod boundary).

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
    if !strings.Contains(resp.Stderr, "alpha_test") {
        t.Fatalf("stderr missing alpha_test:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "beta_test") {
        t.Fatalf("stderr missing beta_test:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "hidden_test") {
        t.Fatalf("stderr should not contain hidden_test (nested go.mod boundary):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "PASS") {
        t.Fatalf("stderr missing PASS:\n%s", resp.Stderr)
    }
}
```
