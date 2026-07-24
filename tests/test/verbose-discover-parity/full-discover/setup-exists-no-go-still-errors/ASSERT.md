---
label: heavy
---

## Expected

- `DiscoverTreeCases` returns a non-empty error.
- Error text includes `must have a Go code block` and identifies the intermediate SETUP path.

## Exit Code

- N/A (API leaf; DiscoverErr is the signal)

```go
import (
	"strings"
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
	if resp.DiscoverErr == "" {
		t.Fatal("expected DiscoverTreeCases error when intermediate SETUP.md exists without Go block")
	}
	low := strings.ToLower(resp.DiscoverErr)
	if !strings.Contains(low, "must have a go code block") {
		t.Fatalf("error must mention Go code block, got:\n%s", resp.DiscoverErr)
	}
	if !strings.Contains(low, "intermediate") || !strings.Contains(low, "setup.md") {
		t.Fatalf("error must identify intermediate SETUP.md path, got:\n%s", resp.DiscoverErr)
	}
}
```
