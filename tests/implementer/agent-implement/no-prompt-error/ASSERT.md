## Expected
- Exit code non-zero.
- Stderr contains "requires".

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode == 0 {
        t.Fatal("expected non-zero exit for missing prompt")
    }
    if !strings.Contains(resp.Stderr, "requires") {
        t.Fatalf("expected 'requires' in stderr, got:\n%s", resp.Stderr)
    }
}
```
