---
label: heavy
---

## Expected

- Exit 0, PASS.
- Plan uses hub: `__hub` and `go test` with suite package.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	out := gotestPlanOut(resp)
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d\n%s", resp.ExitCode, out)
	}
	if !strings.Contains(out, "go test") {
		t.Fatalf("missing go test:\n%s", out)
	}
	// hub dir and suite package (ModeHubSuite)
	if !strings.Contains(out, "__hub") && !strings.Contains(out, "workspace hub") {
		t.Fatalf("expected multi-mod hub plan:\n%s", out)
	}
	if !strings.Contains(out, "PASS") {
		t.Fatalf("missing PASS\n%s", out)
	}
}
```
