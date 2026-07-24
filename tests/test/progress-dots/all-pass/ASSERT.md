---
label: heavy
---

## Expected
- The command succeeds.
- stdout shows 3 dots and summary `(3 Run, 3 Pass, 0 Fail)`.
- stdout MUST NOT contain raw `go test` ok lines like `ok\t...`.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    inline := findInlineSummaryLine(resp.Stdout)
    if inline == "" {
        t.Fatalf("expected inline progress summary with duration, got:\n%s", resp.Stdout)
    }
    if !strings.Contains(inline, "(3 Run, 3 Pass, 0 Fail, 0 Cached) in ") {
        t.Fatalf("expected (3 Run, 3 Pass, 0 Fail, 0 Cached) in <duration>, got %q", inline)
    }
    for _, line := range strings.Split(resp.Stdout, "\n") {
        if strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "ok\t") {
            t.Fatalf("stdout must not contain raw go test ok lines, got: %q\nfull stdout:\n%s", line, resp.Stdout)
        }
    }
}
```
