## Expected

- Non-zero exit code.
- stderr/stdout mentions the unknown flag or usage (error path).
- Must be a real `cache` flag error — not merely `unknown command: cache`
  (that would pass before the feature exists).

## Exit Code

- non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	requireFail(t, resp, err)
	combined := combinedOut(resp)
	if strings.Contains(combined, "unknown command: cache") {
		t.Fatalf("cache not implemented yet (need flag parse after cache is registered):\n%s", combined)
	}
	// Accept common CLI unknown-flag phrasings once cache exists.
	lower := strings.ToLower(combined)
	ok := strings.Contains(combined, "not-a-real-flag") ||
		strings.Contains(lower, "unknown") ||
		strings.Contains(lower, "flag") ||
		strings.Contains(combined, "Usage:")
	if !ok {
		t.Fatalf("expected unknown-flag/usage signal on stderr/stdout:\nstderr=%q\nstdout=%q", resp.Stderr, resp.Stdout)
	}
}
```
