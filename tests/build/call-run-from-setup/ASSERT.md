## Expected
- With `Run := run` alias, descendant Setup code can call `Run(t, d, req)`.
- The build now succeeds (exit 0).

```go
import (
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected doctest test to succeed with Run alias, got exit %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
}
```
