## Expected
- Exit code 0.
- Stdout contains session ID `status-finished-test`.
- Stdout contains status `finished`.
- Stdout contains `opencode` (runner).
- Stdout contains `codex-thread-abc123`.
- Stdout contains `ses_open_xyz789`.
- Stdout contains created_at time.
- Stdout contains `Events: 5 lines`.
- Stdout contains last 3 events (the text event, the task tool_use, and the second text event).
- Stdout contains `Build project` or `Building the project`.
- Stdout contains `All 16 tests pass`.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    stdout := resp.Stdout

    checks := map[string]string{
        "session id":              "status-finished-test",
        "status finished":         "finished",
        "runner":                  "opencode",
        "codex thread id":         "codex-thread-abc123",
        "opencode session id":     "ses_open_xyz789",
        "event count":             "5",
        "event content 1":         "Build project",
        "event content 2":         "Building the project",
        "event content 3":         "All 16 tests pass",
    }

    for label, want := range checks {
        if !strings.Contains(stdout, want) {
            t.Errorf("%s: stdout missing %q, got:\n%s", label, want, stdout)
        }
    }

    if t.Failed() {
        t.Fatalf("full stdout:\n%s\nstderr:\n%s", stdout, resp.Stderr)
    }
}
```
