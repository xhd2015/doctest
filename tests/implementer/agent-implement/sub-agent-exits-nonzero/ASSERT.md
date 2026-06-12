## Expected
- Exit code non-zero.
- Error message appears in stderr or stdout.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode == 0 {
        t.Fatal("expected non-zero exit for sub-agent failure")
    }
    combined := resp.Stderr + resp.Stdout
    if !strings.Contains(combined, "sub-agent") {
        t.Fatalf("expected 'sub-agent' in output:\n%s", combined)
    }
}
```
