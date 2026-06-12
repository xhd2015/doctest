## Expected
- Exit code 0.
- Stdout contains the mock response text.
- Session directory is created.

```go
import (
    "os"
    "path/filepath"
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

    home, _ := os.UserHomeDir()
    agentDir := filepath.Join(home, ".agent-pro", "dedicated-agents", "doctest-agent", "sessions")
    entries, readErr := os.ReadDir(agentDir)
    if readErr != nil {
        t.Fatalf("cannot read session dir %s: %v", agentDir, readErr)
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
