## Expected
- Non-zero exit code.
- Stderr does NOT contain `child_test` (walk down stops at git boundary).
- Stderr contains "no tests found".

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
        t.Fatalf("expected non-zero exit code (no tests found, child repo skipped), got 0\nstderr:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "child_test") {
        t.Fatalf("stderr should NOT contain child_test (separate git repo, walk down should stop at boundary):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "no tests found") {
        t.Fatalf("stderr missing 'no tests found':\n%s", resp.Stderr)
    }
}
```
