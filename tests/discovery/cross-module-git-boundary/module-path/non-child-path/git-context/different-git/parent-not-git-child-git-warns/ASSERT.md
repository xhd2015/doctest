---
label: heavy
---

## Expected
- Child tests are NOT discovered.
- Stderr contains `warning:` with `git repository mismatch`.
- Stderr contains `no tests`.

## Errors
- None from the doctest runner itself; command exits with failure due to no tests.

## Exit Code
- Non-zero

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
        t.Fatalf("expected non-zero exit code (child skipped, no tests), got 0\nstderr:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "warning:") {
        t.Fatalf("stderr missing warning (parent not git, child git should warn):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "git repository mismatch") {
        t.Fatalf("stderr missing 'git repository mismatch' reason:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "parent/cli") {
        t.Fatalf("stderr missing nested module path in warning:\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "no tests") {
        t.Fatalf("stderr missing 'no tests':\n%s", resp.Stderr)
    }
}
```
