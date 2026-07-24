---
label: heavy
---

## Expected
- Soft "no tests" outcome (exit 0 + `no tests` on stderr).
- `./...` from a directory with no DOCTEST.md at or below finds nothing.
- stderr contains "no tests".

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    // Product: ErrNoTestsFound → soft exit 0 + "no tests" on stderr.
    if resp.ExitCode != 0 {
        t.Fatalf("expected soft exit 0 for no tests, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "no tests") {
        t.Fatalf("stderr missing 'no tests':\n%s", resp.Stderr)
    }
}
```
