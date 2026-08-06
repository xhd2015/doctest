## Expected

- `DiscoverTreeCases` succeeds when intermediate SETUP.md is prose-only (no Go block).
- DiscoverErr is empty.

## Exit Code

- N/A (API leaf; DiscoverErr is the signal)

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = req
	if err != nil {
		t.Fatalf("harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.DiscoverErr != "" {
		t.Fatalf("expected no DiscoverTreeCases error for prose-only intermediate SETUP, got:\n%s", resp.DiscoverErr)
	}
}
```
