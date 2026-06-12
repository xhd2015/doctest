## Expected
- Exit code 0.
- The followup is delivered on the same session (resume, not new session).
- The PROMPT template is not included (resume skips template).

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
    if !strings.Contains(resp.Stdout, "resumed and finished") {
        t.Fatalf("stdout missing expected text:\n%s", resp.Stdout)
    }
    if strings.Contains(resp.Stdout, "Step 1: Read the Test Tree") {
        t.Fatalf("stdout should not contain PROMPT template text on resume:\n%s", resp.Stdout)
    }
}
```
