## Expected
- The --agent-runner flag reaches the agent provider.
- The mock config response appears in stdout.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "args forwarded ok") {
        t.Fatalf("stdout missing mock text:\n%s", resp.Stdout)
    }
}
```
