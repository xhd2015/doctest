## Expected
- Exit code non-zero.
- Stderr contains "requires".

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode == 0 {
        t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stderr+resp.Stdout, "requires") {
        t.Fatalf("expected error about missing prompt:\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
    }
}
```
