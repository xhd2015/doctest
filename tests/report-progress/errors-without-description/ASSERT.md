## Expected
- Exit code non-zero.
- Stderr contains usage information mentioning `<description>`.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode == 0 {
        t.Fatal("expected non-zero exit when no description is provided")
    }
    if !strings.Contains(resp.Stderr, "usage") && !strings.Contains(resp.Stderr, "description") {
        t.Fatalf("expected usage message with 'description' in stderr:\n%s", resp.Stderr)
    }
}
```
