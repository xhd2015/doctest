## Expected
- Exit code non-zero.
- Stderr or stdout contains the mock error message.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode == 0 {
        t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stdout+resp.Stderr, "build failed") {
        t.Fatalf("expected 'build failed' in output:\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
    }
}
```
