## Expected
- With `_ = helperName` suppression lines, all helpers are marked as used.
- The build now succeeds (exit 0).

```go
import (
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected doctest test to succeed with helper suppression, got exit %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
}
```
