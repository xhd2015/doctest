---
label: e2e
---

## Expected

- `go run ./suite` exits 0.
- Output contains `from-internal-leaf`.

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
	if !strings.Contains(resp.Combined, "from-internal-leaf") {
		t.Fatalf("expected from-internal-leaf in output:\n%s", resp.Combined)
	}
}
```
