## Expected
- Exit code 0.
- stderr contains "my-feature" (the doctest tree under tests/ was found).
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
        t.Errorf("stderr missing 'my-feature' (doctest tree should be found even without go.mod):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Errorf("stdout missing PASS:\n%s", resp.Stdout)
    }
    // Should NOT report "no tests"
    if strings.Contains(resp.Stderr, "no tests") {
        t.Errorf("stderr contains 'no tests' but should find tests:\n%s", resp.Stderr)
    }
}
```
