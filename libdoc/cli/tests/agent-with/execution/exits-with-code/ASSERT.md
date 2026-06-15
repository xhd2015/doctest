## Expected
- Exit code 42 (propagated from child).

## Exit Code
- 42.

```go
import (
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error (exit 42), got nil")
    }
    if resp.ExitCode != 42 {
        t.Fatalf("expected exit code 42, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
}
```
