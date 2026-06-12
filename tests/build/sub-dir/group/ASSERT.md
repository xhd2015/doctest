## Expected
- The command succeeds.
- Both leaf-1 and leaf-2 are built (2 test cases), leaf-3 is NOT built.

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
    if !strings.Contains(resp.Stderr, "2 test cases") {
        t.Fatalf("expected 2 test cases, stderr:\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "group_b_leaf_3") {
        t.Fatalf("expected leaf-3 NOT to run, stderr:\n%s", resp.Stderr)
    }
}
```
