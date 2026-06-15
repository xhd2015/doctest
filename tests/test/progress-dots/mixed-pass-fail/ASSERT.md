## Expected
- stdout shows 3 dots and summary `(3 Run, 2 Pass, 1 Fail)`.
- Exit code may be non-zero due to failing test.
- stdout MUST NOT contain raw `go test` ok lines like `ok\t...`.
- Failing test lines (`FAIL\t...`) are printed.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if !strings.Contains(resp.Stdout, "...  (3 Run, 2 Pass, 1 Fail, 0 Cached)") {
        t.Fatalf("expected stdout to contain '...  (3 Run, 2 Pass, 1 Fail, 0 Cached)', got:\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
    }
    hasFail := false
    for _, line := range strings.Split(resp.Stdout, "\n") {
        if strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "ok\t") {
            t.Fatalf("stdout must not contain raw go test ok lines, got: %q\nfull stdout:\n%s", line, resp.Stdout)
        }
        if strings.HasPrefix(line, "FAIL\t") {
            hasFail = true
        }
    }
    if !hasFail {
        t.Fatalf("expected stdout to contain failing test lines (FAIL\\t...), got:\n%s", resp.Stdout)
    }
}
```
