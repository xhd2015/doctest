## Expected
- stdout shows 3 dots and summary `(3 Run, 0 Pass, 3 Fail)`.
- Exit code may be non-zero due to failing test.
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
    if !strings.Contains(resp.Stdout, "...  (3 Run, 0 Pass, 3 Fail)") {
        t.Fatalf("expected stdout to contain '...  (3 Run, 0 Pass, 3 Fail)', got:\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
    }
    for _, line := range strings.Split(resp.Stdout, "\n") {
        if strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "ok\t") || strings.HasPrefix(line, "FAIL\t") {
            t.Fatalf("stdout must not contain raw go test output lines, got: %q\nfull stdout:\n%s", line, resp.Stdout)
        }
    }
}
```
