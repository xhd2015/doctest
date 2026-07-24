## Expected
- `resp.Err` is non-nil and contains `"unknown agent runner id"` and `"idonotexist"`.
- Error mentions `"idonotexist"` (from env var override), NOT `"codex"` (from CODEX_THREAD_ID).
- This proves `DOCTEST_SUBAGENT_AGENT_RUNNER` env var takes priority over `CODEX_THREAD_ID`.

## Errors
- `agentprovider.Build` rejects unknown runner IDs.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error from unknown runner, got nil")
    }
    msg := resp.Err.Error()
    if !strings.Contains(msg, "unknown agent runner id") || !strings.Contains(msg, "idonotexist") {
        t.Fatalf("expected error containing 'idonotexist' (env var should beat CODEX_THREAD_ID), got: %v", resp.Err)
    }
    if strings.Contains(msg, "codex") {
        t.Fatalf("expected error to mention 'idonotexist', not 'codex' (env var must beat CODEX_THREAD_ID): %v", resp.Err)
    }
}
```
