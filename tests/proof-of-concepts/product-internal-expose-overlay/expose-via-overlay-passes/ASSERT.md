---
label: e2e
---

## Expected

- `go run -overlay=… ./suite_expose` exits 0.
- Prints `hello-from-app-internal`.
- Virtual expose file is absent under `app/` on disk.

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
	if !strings.Contains(resp.Combined, "hello-from-app-internal") {
		t.Fatalf("expected hello-from-app-internal:\n%s", resp.Combined)
	}
	virt := filepath.Join(req.AppRoot, "__doctest_internal_expose", "greet", "expose.go")
	if _, err := os.Stat(virt); !os.IsNotExist(err) {
		t.Fatalf("expected overlay-only path absent on disk: %s (err=%v)", virt, err)
	}
}
```
