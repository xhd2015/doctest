---
label: heavy
---

## Expected
- doctest test fails with non-zero exit.
- stderr contains an error indicating that the non-root SETUP.md must have
  func Setup.

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
    if !strings.Contains(resp.Stderr, "must have func Setup") {
        t.Fatalf("expected 'must have func Setup' in stderr:\n%s", resp.Stderr)
    }
}
```
