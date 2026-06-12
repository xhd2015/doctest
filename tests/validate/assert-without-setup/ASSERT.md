## Expected
- The command fails with a setup-missing validation error.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode == 0 {
        t.Fatalf("expected nonzero exit")
    }
    if !strings.Contains(resp.Stderr, "ASSERT.md found but SETUP.md missing") {
        t.Fatalf("stderr missing validation error:\n%s", resp.Stderr)
    }
}
```
