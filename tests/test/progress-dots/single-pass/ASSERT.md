## Expected
- The command succeeds.
- stdout shows 1 dot and summary `(1 Run, 1 Pass, 0 Fail)`.
- stdout MUST NOT contain raw `go test` lines like `ok\t...` or `FAIL\t...`.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, ".  (1 Run, 1 Pass, 0 Fail)") {
        t.Fatalf("expected stdout to contain '.  (1 Run, 1 Pass, 0 Fail)', got:\n%s", resp.Stdout)
    }
    for _, line := range strings.Split(resp.Stdout, "\n") {
        if strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "ok\t") || strings.HasPrefix(line, "FAIL\t") {
            t.Fatalf("stdout must not contain raw go test output lines, got: %q\nfull stdout:\n%s", line, resp.Stdout)
        }
    }
}
```
