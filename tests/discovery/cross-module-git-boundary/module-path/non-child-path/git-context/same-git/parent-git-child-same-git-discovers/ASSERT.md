## Expected
- Exit code 0.
- `child_test` is discovered.
- No `warning:` on stderr.
- No `no tests found`.

## Side Effects
- None beyond running nested doctest tree.

## Exit Code
- 0

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
    if strings.Contains(resp.Stderr, "warning:") {
        t.Fatalf("stderr should not contain warning (same git repo, should discover):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "child_test") {
        t.Fatalf("stderr missing child_test (same git repo should discover sibling module):\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "no tests found") {
        t.Fatalf("stderr should not contain 'no tests found':\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "PASS") {
        t.Fatalf("stdout missing PASS:\n%s", resp.Stdout)
    }
}
```
