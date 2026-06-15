## Expected
- Exit code 0.
- Stdout contains `arg1` and `arg2`.

## Exit Code
- 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.Err != nil {
        t.Fatalf("unexpected error: %v", resp.Err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "arg1") || !strings.Contains(resp.Stdout, "arg2") {
        t.Fatalf("expected stdout to contain 'arg1' and 'arg2', got:\n%s", resp.Stdout)
    }
}
```
