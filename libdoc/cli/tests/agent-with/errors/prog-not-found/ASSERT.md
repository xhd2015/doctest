## Expected
- `resp.Err` is non-nil and contains `"executable file not found"` or `"no such file"`.

## Errors
- The subcommand must return an error when the program is not found in PATH.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error for nonexistent prog, got nil")
    }
    msg := resp.Err.Error()
    if !strings.Contains(msg, "executable file not found") && !strings.Contains(msg, "no such file") {
        t.Fatalf("expected error about missing executable, got: %v", resp.Err)
    }
}
```
