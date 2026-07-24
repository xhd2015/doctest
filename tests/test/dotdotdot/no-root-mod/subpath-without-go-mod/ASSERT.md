---
label: heavy
---

## Expected
- Exit code 0.
- stderr contains "my-feature" (the doctest tree under tests/ was found).
- stderr does NOT contain "other-feature" (the tree under other/ was excluded by the path prefix).
- stderr contains PASS.

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
        t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "my-feature") {
        t.Errorf("stderr missing 'my-feature' (doctest tree under tests/ should be found):\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "other-feature") {
        t.Errorf("stderr should not contain 'other-feature' (excluded by tests/... prefix):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Errorf("stdout missing PASS:\n%s", resp.Stdout)
    }
    // Should NOT report "no tests" or "file does not exist"
    if strings.Contains(resp.Stderr, "no tests") {
        t.Errorf("stderr contains 'no tests' but should find tests:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "file does not exist") {
        t.Errorf("stderr contains 'file does not exist' but the path exists:\n%s", resp.Stderr)
    }
}
```
