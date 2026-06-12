## Expected
- Exit code non-zero (no session ID can be resolved).
- Stderr contains "cannot detect session id".
- Stderr contains a generated session ID (e.g. `gen_` prefix).
- Stderr mentions `--session-id` flag usage.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode == 0 {
        t.Fatal("expected non-zero exit code when no session ID source is available")
    }

    combined := resp.Stdout + resp.Stderr

    if !strings.Contains(combined, "cannot detect session id") {
        t.Fatalf("expected 'cannot detect session id' in output:\n%s", combined)
    }
    if !strings.Contains(combined, "--session-id") {
        t.Fatalf("expected '--session-id' flag in error message:\n%s", combined)
    }
}
```
