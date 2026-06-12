## Expected
- With topological sort of types, the build now succeeds (exit 0).
- LocationEntry is emitted after GitInfo, resolving the forward reference.

```go
import (
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected doctest test to succeed with sorted types, got exit %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
}
```
