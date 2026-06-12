## Expected
- Exit code non-zero (tests fail).
- Output contains "FAIL" and "not implemented".

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode == 0 {
        t.Fatalf("expected non-zero exit (RED), got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stdout+resp.Stderr, "not implemented") {
        t.Fatalf("expected output to contain 'not implemented':\n%s\n%s", resp.Stdout, resp.Stderr)
    }
}
```
