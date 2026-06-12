## Expected
- The command succeeds using mod-b's SETUP.md.
- Only 1 test case is built.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v\nstderr:\n%s", err, resp.Stderr)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected zero exit, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "1 test case") {
        t.Fatalf("expected 1 test case, stderr:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "leaf_a") {
        t.Fatalf("expected leaf-a (mod-a) NOT to run, stderr:\n%s", resp.Stderr)
    }
}
```
