---
label: e2e
---

## Expected

- `go list -overlay=…` exits 0.
- Output contains the virtual package import path.

## Exit Code

- Zero.

```go
import (
	"strings"
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
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\n%s", resp.ExitCode, resp.Combined)
	}
	want := "example.com/realmod/http/__doctest_internal_shim_overlay/leaf"
	if !strings.Contains(resp.Combined, want) {
		t.Fatalf("expected %q in output:\n%s", want, resp.Combined)
	}
}
```
