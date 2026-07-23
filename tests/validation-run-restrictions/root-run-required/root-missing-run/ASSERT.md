---
label: heavy
---

## Expected
- doctest test fails with non-zero exit.
- stderr contains the error message indicating "must have func Run".

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
        t.Fatalf("expected nonzero exit, stdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "must have func Run") && !strings.Contains(resp.Stdout, "must have func Run") {
        t.Fatalf("expected 'must have func Run' in output, stdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
    }
}
```
