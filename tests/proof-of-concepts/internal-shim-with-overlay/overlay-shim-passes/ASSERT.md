---
label: e2e
---

## Expected

- `go run -overlay=overlay.json ./suite_overlay` exits 0.
- Output contains `from-internal-leaf`.
- Virtual shim path is not a real directory under the workdir (only overlay).

## Exit Code

- Zero.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
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
	// Virtual path must not exist on disk (overlay-only).
	virt := filepath.Join(req.WorkDir, "http", "__doctest_internal_shim_overlay", "leaf", "shim.go")
	if _, err := os.Stat(virt); !os.IsNotExist(err) {
		t.Fatalf("expected overlay-only path absent on disk: %s (err=%v)", virt, err)
	}
}
```
