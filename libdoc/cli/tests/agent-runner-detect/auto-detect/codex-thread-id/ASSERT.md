## Expected
- `resp.Err` is non-nil.
- Error mentions `"codex"` (the runner selected by `CODEX_THREAD_ID` detection).
- This proves `CODEX_THREAD_ID` triggers codex detection in the subagent.

## Errors
- `codex` binary not found in test environment (expected: `"codex not found"` or similar).

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error (codex not found), got nil")
    }
    msg := resp.Err.Error()
    if !strings.Contains(msg, "codex") {
        t.Fatalf("expected error to mention 'codex' (from CODEX_THREAD_ID detection), got: %v", resp.Err)
    }
}
```
