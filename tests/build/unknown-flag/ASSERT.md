## Expected
- The command fails.
- stderr reports the unknown runner option.

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
    if !strings.Contains(resp.Stderr, "unknown flag") &&
        !strings.Contains(resp.Stderr, "unknown option") &&
        !strings.Contains(resp.Stderr, "unrecognized flag") {
        t.Fatalf("stderr missing unknown flag message:\n%s", resp.Stderr)
    }
}
```
