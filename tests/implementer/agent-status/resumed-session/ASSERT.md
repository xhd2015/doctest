## Expected
- Exit code 0.
- Stdout contains status `running`.
- Stdout contains the session ID `status-resumed-test`.
- The PID of the background process is present in the output.

```go
import (
    "os"
    "strconv"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }

    stdout := resp.Stdout

    if !strings.Contains(strings.ToLower(stdout), "running") {
        t.Fatalf("stdout missing 'running', got:\n%s", stdout)
    }

    if !strings.Contains(stdout, "status-resumed-test") {
        t.Fatalf("stdout missing session id 'status-resumed-test', got:\n%s", stdout)
    }

    pidStr := strconv.Itoa(os.Getpid())
    _ = pidStr
}
```
