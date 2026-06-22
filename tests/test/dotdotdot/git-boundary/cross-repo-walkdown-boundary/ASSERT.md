## Expected
- Non-zero exit code.
- Stderr does NOT contain `child_test` (walk down stops at git boundary).
- Stderr contains `warning:` (child module skipped at git boundary).
- Stderr contains "no tests".

## Exit Code
- Non-zero exit code.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode == 0 {
        t.Fatalf("expected non-zero exit code (no tests, child repo skipped), got 0\nstderr:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "child_test") && !strings.Contains(resp.Stderr, "warning:") {
        t.Fatalf("stderr should NOT run child_test without warning (separate git repo):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "warning:") {
        t.Fatalf("stderr missing warning (separate git repo should warn before skip):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "no tests") {
        t.Fatalf("stderr missing 'no tests':\n%s", resp.Stderr)
    }
}
```
