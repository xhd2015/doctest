## Expected
- Exit code 0.
- Stdout contains the sub-agent's response text.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "I have implemented the feature.") {
        t.Fatalf("stdout missing expected text:\n%s", resp.Stdout)
    }
}
```
