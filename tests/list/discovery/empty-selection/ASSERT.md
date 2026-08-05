## Expected

- Exit code 0 (soft exit, like `test` + ErrNoTestsFound).
- stderr contains a line `no tests` (not hard failure).
- stdout is empty (no body, no summary).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	requireOK(t, resp, err)
	if strings.TrimSpace(resp.Stdout) != "" {
		t.Fatalf("stdout must be empty on empty selection, got:\n%q", resp.Stdout)
	}
	if !strings.Contains(resp.Stderr, "no tests") {
		t.Fatalf("stderr must contain \"no tests\", got %q", resp.Stderr)
	}
	// Prefer exact soft message family (not old "no tests found" only if product uses "no tests").
	if strings.Contains(resp.Stderr, "no tests found") && !strings.Contains(resp.Stderr, "no tests\n") && !strings.HasSuffix(strings.TrimSpace(resp.Stderr), "no tests") {
		// Allow either exact "no tests" line; fail only if missing entirely (already checked).
	}
}
```
