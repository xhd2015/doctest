## Expected
- Exit code 0.
- `./...` from `group/subgroup/tests/` only finds that tree.
- Stderr contains `deep_tests`, does NOT contain `other`, `alpha_test`, or `beta_test`.

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
    if !strings.Contains(resp.Stderr, "deep_tests") {
        t.Fatalf("stderr missing deep_tests:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "other") {
        t.Fatalf("stderr should not contain other (sibling, not under working dir):\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "alpha_test") {
        t.Fatalf("stderr should not contain alpha_test (not under working dir):\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "beta_test") {
        t.Fatalf("stderr should not contain beta_test (not under working dir):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Fatalf("stdout missing PASS:\n%s", resp.Stdout)
    }
}
```
