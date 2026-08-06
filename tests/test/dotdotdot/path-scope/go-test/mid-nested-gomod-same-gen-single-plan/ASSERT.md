## Expected

- **Never** two `go test` lines with the same `cd` dir (must combine).
- Shared gen: exactly **one** plan line; multi-pattern includes nested suite path and mid path.
- Markers once each; no sibling.

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
	// Horn 1 of the rule: same Dir must not get multiple go test processes.
	pathScopeAssertPlansByDir(t, plans, out)

	// Shared --gen-dir: all plans share one dir → exactly one combined process.
	if len(plans) != 1 {
		t.Fatalf("S4 shared-gen want exactly 1 go test (combine patterns), got %d:\n  %s\nfull:\n%s",
			len(plans), strings.Join(plans, "\n  "), out)
	}
	p := plans[0]
	// Multi-pattern: both nested-mod suite packages and mid packages.
	hasSuite := strings.Contains(p, "./suite") || strings.Contains(p, "suite/...")
	hasMid := strings.Contains(p, "mid") && strings.Contains(p, "/...")
	if !hasSuite || !hasMid {
		t.Fatalf("combined plan want both suite scope and mid path ...:\n  %s\nfull:\n%s", p, out)
	}
	if n := pathScopeCountMarker(out, "MID_LEAF"); n != 1 {
		t.Fatalf("MARKER:MID_LEAF want 1 got %d\n%s", n, out)
	}
	if n := pathScopeCountMarker(out, "NESTED_MOD_LEAF"); n != 1 {
		t.Fatalf("MARKER:NESTED_MOD_LEAF want 1 got %d\n%s", n, out)
	}
	if pathScopeCountMarker(out, "SIBLING_LEAF") != 0 {
		t.Fatalf("sibling must not run\n%s", out)
	}
}
```
