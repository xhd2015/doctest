## Expected
- Exit code 0.
- Stdout contains formatted trace output with the event text.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    if !strings.Contains(resp.Stdout, "Hello from trace!") {
        t.Fatalf("stdout missing event text 'Hello from trace!', got:\n%s", resp.Stdout)
    }
}
```
