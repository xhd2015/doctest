## Expected
- The command exits 0 with `no tests` on stderr.

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
        t.Fatalf("expected zero exit, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "no tests") {
        t.Fatalf("stderr missing 'no tests':\n%s", resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "no tests found") {
        t.Fatalf("stderr must not contain old 'no tests found' message:\n%s", resp.Stderr)
    }
}
```