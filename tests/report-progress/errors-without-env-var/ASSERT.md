## Expected
- Exit code non-zero.
- Stderr contains "PROGRESS_FILE must be set".

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode == 0 {
        t.Fatal("expected non-zero exit when PROGRESS_FILE is not set")
    }
    if !strings.Contains(resp.Stderr, "PROGRESS_FILE must be set") {
        t.Fatalf("expected 'PROGRESS_FILE must be set' in stderr:\n%s", resp.Stderr)
    }
}
```
