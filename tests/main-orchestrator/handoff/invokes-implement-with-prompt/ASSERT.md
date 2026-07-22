---
label: heavy
---

## Expected
- Exit code 0.
- Stdout contains the mock response text.
- Session directory is created.

```go
import (
    "os"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "implemented greet feature") {
        t.Fatalf("stdout missing mock text:\n%s", resp.Stdout)
    }

    sessionHome := req.SessionHome
    if sessionHome == "" {
        t.Fatal("DOCTEST_DEBUG_SESSION_HOME not set")
    }
    entries, readErr := os.ReadDir(sessionHome)
    if readErr != nil {
        t.Fatalf("cannot read session dir %s: %v", sessionHome, readErr)
    }
    found := false
    for _, entry := range entries {
        if entry.IsDir() {
            found = true
            break
        }
    }
    if !found {
        t.Fatal("no session directory created")
    }
}
```
