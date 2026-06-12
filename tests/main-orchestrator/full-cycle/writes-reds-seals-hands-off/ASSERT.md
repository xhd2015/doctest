## Expected
- All phases succeed: RED confirmed, git staged, sub-agent completed.
- Exit code 0.
- Stdout contains the sub-agent completion text.

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
    if !strings.Contains(resp.Stdout, "feature implemented") {
        t.Fatalf("stdout missing completion text:\n%s", resp.Stdout)
    }

    // Verify session was created
    home, _ := os.UserHomeDir()
    agentDir := filepath.Join(home, ".agent-pro", "dedicated-agents", "doctest-agent", "sessions")
    entries, _ := os.ReadDir(agentDir)
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
