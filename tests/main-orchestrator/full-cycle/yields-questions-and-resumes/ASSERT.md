## Expected
- Exit code 0 (yielding questions is not an error).
- Stdout contains the question text.
- Session exists for resume.

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
    if !strings.Contains(resp.Stdout, `"question"`) {
        t.Fatalf("stdout missing question:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stdout, `What should Greet`) {
        t.Fatalf("stdout missing question text:\n%s", resp.Stdout)
    }

    sessionHome := os.Getenv("DOCTEST_DEBUG_SESSION_HOME")
    if sessionHome == "" {
        t.Fatal("DOCTEST_DEBUG_SESSION_HOME not set")
    }
    entries, _ := os.ReadDir(sessionHome)
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
