## Expected
- The `validate` command is no longer recognised.
- stderr reports "unknown command: validate".

## Exit Code
- A nonzero exit code.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode == 0 {
        t.Fatalf("expected nonzero exit, stdout:\n%s", resp.Stdout)
    }
    if !strings.Contains(resp.Stderr, "unknown command: validate") {
        t.Fatalf("stderr missing unknown-command message:\n%s", resp.Stderr)
    }
}
```
