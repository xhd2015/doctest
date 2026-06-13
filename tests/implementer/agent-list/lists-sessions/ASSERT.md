## Expected
- Exit code 0.
- Stdout contains all three session IDs: `list-alpha`, `list-beta`, `list-gamma`.
- Stdout contains both runners: `opencode` and `codex`.
- Stdout contains time values.

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
        "session alpha":  "list-alpha",
        "session beta":   "list-beta",
        "session gamma":  "list-gamma",
        "runner opencode":"opencode",
        "runner codex":   "codex",
    }

    for label, want := range checks {
        if !strings.Contains(stdout, want) {
            t.Errorf("%s: stdout missing %q", label, want)
        }
    }

    if t.Failed() {
        t.Fatalf("full stdout:\n%s\nstderr:\n%s", stdout, resp.Stderr)
    }
}
```
