## Expected
- Non-zero exit code.
- Stderr contains error message about bare `...` not being supported.
- The message should direct users to use `./...` or `path/...` instead.

## Exit Code
- Non-zero exit code.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp == nil {
        t.Fatal("resp is nil")
    }
    if resp.ExitCode == 0 {
        t.Fatalf("expected non-zero exit code for bare '...', got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
    }
    combined := resp.Stdout + resp.Stderr
    if !strings.Contains(combined, "...") {
        t.Fatalf("expected error message mentioning '...', got:\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
    }
}
```
