---
label: heavy
---

## Expected
- Exit code 0.
- `./...` from `alpha_test/` only finds doctest trees at or below the working directory.
- Stderr contains `alpha_test`, does NOT contain `beta_test` or `hidden_test`.

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
        t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "alpha_test") {
        t.Fatalf("stderr missing alpha_test:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "beta_test") {
        t.Fatalf("stderr should not contain beta_test (sibling, not under working dir):\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "hidden_test") {
        t.Fatalf("stderr should not contain hidden_test (nested go.mod boundary):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Fatalf("stdout missing PASS:\n%s", resp.Stdout)
    }
}
```
