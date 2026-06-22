## Expected
- Non-zero exit code.
- Stderr does NOT contain `parent_test` (parent is above CWD, `./...` only looks down).
- Stderr contains "no tests" (no tests in or below CWD).

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
        t.Fatalf("expected non-zero exit code (no tests in or below CWD), got 0\nstderr:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "parent_test") {
        t.Fatalf("stderr should NOT contain parent_test (./... only looks down from CWD):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "no tests") {
        t.Fatalf("stderr missing 'no tests':\n%s", resp.Stderr)
    }
}
```
