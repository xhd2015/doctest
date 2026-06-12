## Expected
- The command succeeds.
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
    if strings.Contains(resp.Stderr, "group_a_leaf_2") {
        t.Fatalf("expected leaf-2 NOT to run, stderr:\n%s", resp.Stderr)
    }
}
```
