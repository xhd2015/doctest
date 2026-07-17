---
label: heavy
---

## Expected
- Exit code 0.
- Stdout contains `from --session-id, matching CODEX_THREAD_ID` (fallback match indicator).
- Stdout contains session ID `codex-fallback-test`.
- Stdout contains status `finished`.
- Stdout contains `opencode` (runner).
- Stdout contains `Events: 2 lines`.

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

    if !strings.Contains(stdout, "from --session-id, matching CODEX_THREAD_ID") {
        t.Errorf("stdout missing fallback match indicator, got:\n%s", stdout)
    }

    checks := map[string]string{
        "session id":        "codex-fallback-test",
        "status finished":   "finished",
        "runner":            "opencode",
        "event count":       "2",
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
