## Expected
- `resp.Err` is non-nil and contains `"--agent-runner requires a value"`.

## Errors
- The parser must reject an empty `--agent-runner` value.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error for missing --agent-runner value, got nil")
    }
    if !strings.Contains(resp.Err.Error(), "--agent-runner requires a value") {
        t.Fatalf("expected error '--agent-runner requires a value', got: %v", resp.Err)
    }
}
```
