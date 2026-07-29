---
label: heavy
---

## Expected

- Exit 0, PASS.
- **Exactly one** `cd … && go test …` plan line for the workspace suite family
  (`__workspace/suite`); no multi-cmd fiction / no hub.

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
	// Single-cmd ModeWorkspaceSuite contract (Phase 1 production path).
	assertExactlyOneGoTestPlanFamily(t, out, "__workspace/suite")
	if strings.Contains(out, "__hub") {
		t.Fatalf("single-mod should not use hub:\n%s", out)
	}
	if !strings.Contains(out, "PASS") {
		t.Fatalf("missing PASS\n%s", out)
	}
}
```
