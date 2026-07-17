---
label: heavy
---

## Expected
- Exit code 0.
- Stdout contains status `running`.
- Stdout contains the current PID.
- Stdout contains session ID `status-running-test`.

```go
import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    stdout := resp.Stdout
    pidStr := strconv.Itoa(os.Getpid())

    if !strings.Contains(strings.ToLower(stdout), "running") {
        t.Fatalf("stdout missing 'running', got:\n%s", stdout)
    }

    if !strings.Contains(stdout, pidStr) {
        t.Fatalf("stdout missing PID %s, got:\n%s", pidStr, stdout)
    }

    if !strings.Contains(stdout, "status-running-test") {
        t.Fatalf("stdout missing session id 'status-running-test', got:\n%s", stdout)
    }

    fmt.Printf("stdout:\n%s\n", stdout)
}
```
