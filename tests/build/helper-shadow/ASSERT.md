## Expected
- The validation rule now catches helper redefinition before Go compilation.
- The error message indicates the helper name collision.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed unexpectedly: %v", err)
    }
    if resp.ExitCode == 0 {
        t.Fatal("expected doctest test to fail with helper redefinition validation error, but got exit 0")
    }
    if !strings.Contains(resp.Stderr, "already defined by an ancestor") {
        t.Fatalf("expected helper redefinition validation error in stderr, got:\n%s", resp.Stderr)
    }
}
```
