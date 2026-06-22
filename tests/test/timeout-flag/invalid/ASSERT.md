## Expected
- The command rejects the invalid `--timeout` value.
- Exit code is non-zero.
- stderr mentions `--timeout`.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    combined := resp.Stdout + resp.Stderr
    lower := strings.ToLower(combined)

    if !strings.Contains(lower, "--timeout") {
        t.Fatalf("expected --timeout in error output, got:\n%s", combined)
    }

    if resp.ExitCode == 0 {
        t.Fatalf("expected non-zero exit for invalid --timeout, got 0\ncombined:\n%s", combined)
    }
}
```