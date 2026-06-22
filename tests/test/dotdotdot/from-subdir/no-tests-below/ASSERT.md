## Expected
- Nonzero exit code.
- `./...` from a directory with no DOCTEST.md at or below finds nothing.
- stderr contains "no tests".

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode == 0 {
        t.Fatalf("expected nonzero exit, stderr:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "no tests") {
        t.Fatalf("stderr missing 'no tests':\n%s", resp.Stderr)
    }
}
```
