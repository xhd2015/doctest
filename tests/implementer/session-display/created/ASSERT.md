## Expected
- The command succeeds.
- stdout contains "Session created" message with the session ID.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "Session created: sess-display-create-1") {
        t.Fatalf("expected stdout to contain 'Session created: sess-display-create-1', got:\n%s", resp.Stdout)
    }
}
```
