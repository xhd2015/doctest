## Expected

- Non-zero exit code.
- stderr mentions the unknown flag or usage (error path).
- Must be a real `list` flag error — not merely `unknown command: list`
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
	combined := resp.Stderr + "\n" + resp.Stdout
	if strings.Contains(combined, "unknown command: list") {
		t.Fatalf("list not implemented yet (need flag parse after list is registered):\n%s", combined)
	}
	// Accept common CLI unknown-flag phrasings once list exists.
	ok := strings.Contains(combined, "not-a-real-flag") ||
		strings.Contains(strings.ToLower(combined), "unknown") ||
		strings.Contains(strings.ToLower(combined), "flag") ||
		strings.Contains(combined, "Usage:")
	if !ok {
		t.Fatalf("expected unknown-flag/usage signal on stderr/stdout:\nstderr=%q\nstdout=%q", resp.Stderr, resp.Stdout)
	}
}
```
