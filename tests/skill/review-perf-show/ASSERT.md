## Expected
- The command succeeds.
- stdout includes the doctest-review-perf skill name, description cues, performance budget, workflow practices, and contrast with design review.

## Exit Code
- Exit code 0.

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
        t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stdout+"\n"+resp.Stderr)
    }
    for _, want := range []string{
        "doctest-review-perf",
        "default-suite performance",
        "3 minutes",
        "metrics top",
        "--unlabeled-only",
        "--label-all",
        "session cache",
        "doctest-review",
        "WARNING",
        "perf report",
    } {
        if !strings.Contains(resp.Stdout, want) {
            t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
        }
    }
}
```
