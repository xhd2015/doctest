## Expected
- The command succeeds.
- Count and verbose flags are forwarded to the underlying runner.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected zero exit, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "TestDoctestSuite") {
        t.Fatalf("expected verbose output to include TestDoctestSuite, stdout:\n%s", resp.Stdout)
    }
}
```
