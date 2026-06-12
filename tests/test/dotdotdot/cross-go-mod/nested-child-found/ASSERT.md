## Expected
- Exit code 0.
- `parent_test` is discovered (in root module).
- `child_test` IS discovered (nested module IS a child path).

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
    if !strings.Contains(resp.Stderr, "parent_test") {
        t.Fatalf("stderr missing parent_test:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "child_test") {
        t.Fatalf("stderr missing child_test (nested module IS a child path):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "PASS") {
        t.Fatalf("stderr missing PASS:\n%s", resp.Stderr)
    }
}
```
