## Expected
- Non-zero exit (no session ID source available).
- Stderr mentions `--session-id` flag.
- Stderr mentions `DOCTEST_AGENT_IMPLEMENTER_SESSION_ID` env var.
- Stderr mentions `CODEX_THREAD_ID` env var.

## Exit Code
- Non-zero when no session ID source is available.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode == 0 {
        t.Fatal("expected non-zero exit when no session ID source is available")
    }
    combined := resp.Stdout + resp.Stderr
    if !strings.Contains(combined, "cannot detect session id") {
        t.Fatalf("expected 'cannot detect session id' in output:\n%s", combined)
    }
    if !strings.Contains(combined, "--session-id") {
        t.Fatalf("expected '--session-id' flag in error:\n%s", combined)
    }
}
```
