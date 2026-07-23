---
label: heavy
---

## Expected
- Soft "no tests" outcome (exit 0 + `no tests` on stderr — not a hard fail).
- Stderr does NOT contain `child_test` as a run target (walk down stops at git boundary).
- Stderr contains `warning:` (child module skipped at git boundary).
- Stderr contains "no tests".

## Exit Code
- 0 (soft no-tests after skip warning).

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    // Product: empty discovery after git-boundary skip is soft-exit 0 + "no tests".
    if resp.ExitCode != 0 {
        t.Fatalf("expected soft exit 0 for no tests (child repo skipped), got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }
    if strings.Contains(resp.Stderr, "child_test") && !strings.Contains(resp.Stderr, "warning:") {
        t.Fatalf("stderr should NOT run child_test without warning (separate git repo):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "warning:") {
        t.Fatalf("stderr missing warning (separate git repo should warn before skip):\n%s", resp.Stderr)
    }
    if !strings.Contains(resp.Stderr, "no tests") {
        t.Fatalf("stderr missing 'no tests':\n%s", resp.Stderr)
    }
}
```

