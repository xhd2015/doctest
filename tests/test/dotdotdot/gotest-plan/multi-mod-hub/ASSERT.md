## Expected

- Exit 0, PASS.
- **Exactly one** `cd … && go test …` plan line for ModeHubSuite (`__hub` /
  workspace hub + suite package). No multi-cmd merge path.

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
	// Single-cmd ModeHubSuite contract (Phase 1 production path).
	n := countCdGoTestPlanLines(out)
	if n != 1 {
		t.Fatalf("want exactly 1 cd…&& go test plan line (hub single-cmd), got %d:\n%s", n, out)
	}
	if !strings.Contains(out, "go test") {
		t.Fatalf("missing go test:\n%s", out)
	}
	// hub dir and suite package (ModeHubSuite)
	if !strings.Contains(out, "__hub") && !strings.Contains(out, "workspace hub") {
		t.Fatalf("expected multi-mod hub plan:\n%s", out)
	}
	// Prefer explicit suite package on the hub plan line when present.
	if strings.Contains(out, "__hub") && !strings.Contains(out, "./suite") {
		t.Fatalf("hub plan should run ./suite:\n%s", out)
	}
	if !strings.Contains(out, "PASS") {
		t.Fatalf("missing PASS\n%s", out)
	}
}
```
