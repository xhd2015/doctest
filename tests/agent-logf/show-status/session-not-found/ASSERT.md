## Expected
- stderr contains "session not found" error message.
- stderr message does NOT have a timestamp prefix `[...]`.
- Stdout is empty (exit code 0 since errors go to stderr).

```go
import (
    "regexp"
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }

    if !strings.Contains(resp.Stderr, "session not found") {
        t.Fatalf("expected 'session not found' in stderr, got:\n%q", resp.Stderr)
    }

    tsPrefix := regexp.MustCompile(`\[\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\]`)
    if tsPrefix.MatchString(resp.Stderr) {
        t.Fatalf("stderr error line has unexpected timestamp prefix:\n%q", resp.Stderr)
    }
}
```
