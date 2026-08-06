## Expected

- Exactly **one** `cd … && go test` plan line.
- Pattern is path-scoped under mid: contains `./tree/mid/...` (or equivalent mid `...`).
- No separate nested-only plan (`./tree/mid/nested/...` alone as a second same-dir cmd).
- `MARKER:MID_LEAF` and `MARKER:NESTED_LEAF` each exactly once; no sibling.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	out := pathScopeOut(resp)
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d\n%s", resp.ExitCode, out)
	}
	plans := pathScopeGoTestPlanLines(out)
	pathScopeAssertPlansByDir(t, plans, out)
	if len(plans) != 1 {
		t.Fatalf("S1 want exactly 1 go test plan (parent ./tree/mid/... covers nested), got %d:\n  %s\nfull:\n%s",
			len(plans), strings.Join(plans, "\n  "), out)
	}
	p := plans[0]
	if !strings.Contains(p, "mid") || !strings.Contains(p, "/...") {
		t.Fatalf("plan want mid path ...:\n  %s\n%s", p, out)
	}
	// Parent covering pattern — not nested-only as the sole package arg.
	if strings.Contains(p, "./tree/mid/nested/...") && !strings.Contains(p, "./tree/mid/...") {
		t.Fatalf("plan is nested-only without parent mid/...:\n  %s", p)
	}
	if !strings.Contains(p, "./tree/mid/...") {
		t.Fatalf("plan must include ./tree/mid/... (covers nested):\n  %s\n%s", p, out)
	}
	if n := pathScopeCountMarker(out, "MID_LEAF"); n != 1 {
		t.Fatalf("MARKER:MID_LEAF want 1 got %d\n%s", n, out)
	}
	if n := pathScopeCountMarker(out, "NESTED_LEAF"); n != 1 {
		t.Fatalf("MARKER:NESTED_LEAF want 1 got %d (double-run = redundant nested job)\n%s", n, out)
	}
	if pathScopeCountMarker(out, "SIBLING_LEAF") != 0 {
		t.Fatalf("sibling must not run\n%s", out)
	}
}
```
