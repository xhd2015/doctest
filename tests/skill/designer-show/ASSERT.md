## Expected
- The command succeeds.
- stdout includes the designer prompt content.
- stdout does not mention legacy progress/question CLIs.

## Exit Code
- Exit code 0.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    for _, want := range []string{"Designer", "doctest test", "Step 1: Understand the Requirement", "Questions"} {
        if !strings.Contains(resp.Stdout, want) {
            t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
        }
    }
    for _, absent := range []string{"report-progress", "yield-pending-questions"} {
        if strings.Contains(resp.Stdout, absent) {
            t.Fatalf("stdout must not contain %q:\n%s", absent, resp.Stdout)
        }
    }
}
```
