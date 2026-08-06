---
label: e2e
---

## Expected

- `go build ./suite_direct` fails.
- Output mentions product internal path.

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
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit\n%s", resp.Combined)
	}
	if !containsInternalDenied(resp.Combined) {
		t.Fatalf("expected internal denied:\n%s", resp.Combined)
	}
	if !strings.Contains(resp.Combined, "example.com/app/internal") {
		t.Fatalf("expected example.com/app/internal in:\n%s", resp.Combined)
	}
}
```
