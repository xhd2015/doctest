---
label: heavy
---

## Expected
- Soft "no tests" outcome (product exits 0 and prints `no tests` — not a hard error).
- Stderr does NOT contain `parent_test` (parent is above CWD, `./...` only looks down).
- Stderr contains "no tests" (no tests in or below CWD).

## Exit Code
- 0 (soft no-tests), with `no tests` on stderr.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    // Product: ErrNoTestsFound → exit 0 + "no tests" on stderr (not a hard fail).
    if resp.ExitCode != 0 {
        t.Fatalf("expected soft exit 0 for no tests, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "parent_test") {
        t.Fatalf("stderr should NOT contain parent_test (./... only looks down from CWD):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "no tests") {
        t.Fatalf("stderr missing 'no tests':\n%s", resp.Stderr)
    }
}
```
