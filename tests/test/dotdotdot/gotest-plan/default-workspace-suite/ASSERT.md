---
label: heavy
---

## Expected

- Exit 0, PASS.
- Verbose plan contains `go test` and `__workspace/suite` (ModeWorkspaceSuite).

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
	assertContainsGoTestLine(t, out, "__workspace/suite")
	if strings.Contains(out, "__hub") {
		t.Fatalf("single-mod should not use hub:\n%s", out)
	}
	if !strings.Contains(out, "PASS") {
		t.Fatalf("missing PASS\n%s", out)
	}
}
```
