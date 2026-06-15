## Expected
- `resp.Err` is non-nil and contains `"unknown agent runner id"` and `"idonotexist"`.
- The runner ID passed via `--agent-runner` is used as-is, without auto-detection.

## Errors
- `agentprovider.Build` rejects unknown runner IDs.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error for unknown runner ID, got nil")
    }
    msg := resp.Err.Error()
    if !strings.Contains(msg, "unknown agent runner id") || !strings.Contains(msg, "idonotexist") {
        t.Fatalf("expected error containing 'unknown agent runner id' and 'idonotexist', got: %v", resp.Err)
    }
}
```
