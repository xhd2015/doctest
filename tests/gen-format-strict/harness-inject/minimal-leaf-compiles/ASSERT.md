---
label: heavy
---

## Expected

- Suite generate + `go test` succeed (`resp.RunErr` empty).
- Generated leaf may contain harness session import even though author SETUP did not import session.
- Author-facing rule: user only needed `import "testing"` for their own symbols.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.RunErr != "" {
		t.Fatalf("A4: expected minimal harness-inject leaf to compile, RunErr=%v\nstdout:\n%s\nstderr:\n%s\nleaf:\n%s",
			resp.RunErr, resp.Stdout, resp.Stderr, resp.LeafGo)
	}
	// Optional signal: harness inject still present when leaf source captured.
	if resp.LeafGo != "" && !containsSessionImport(resp.LeafGo) {
		// Some layouts may put inject only in shared packages; do not hard-fail.
		t.Logf("note: leaf Go lacked session import; suite still passed (ok if inject elsewhere)\nleaf:\n%s", resp.LeafGo)
	}
}
```
