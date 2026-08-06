---
label: e2e
---

## Expected

- `go build ./suite_direct` fails.
- Combined output mentions internal package not allowed.

## Exit Code

- Non-zero.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0\n%s", resp.Combined)
	}
	if !containsInternalDenied(resp.Combined) {
		t.Fatalf("expected internal package denied message, got:\n%s", resp.Combined)
	}
}
```
