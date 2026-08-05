## Expected

- Exit 0.
- Exactly one body line for the root: leaves=2, L2:L3=2:0 (100.0%/0.0%), unlabeled=2.
- Summary present: roots=1, leaves=2, matching L2:L3 and labels.
- Body precedes summary (`---` separator).
- Trailing newline after last summary line.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	requireOK(t, resp, err)
	rep := parseListReport(t, resp.Stdout)
	if len(rep.Body) != 1 {
		t.Fatalf("want 1 body line, got %d\n%s", len(rep.Body), resp.Stdout)
	}
	if !rep.HasSep {
		t.Fatalf("expected blank+--- summary separator:\n%s", resp.Stdout)
	}
	b := findBody(t, rep.Body, req.Roots[0])
	if b.Leaves != 2 || b.L2 != 2 || b.L3 != 0 {
		t.Fatalf("body stats = leaves=%d L2:L3=%d:%d want 2 / 2:0", b.Leaves, b.L2, b.L3)
	}
	if b.P2 != "100.0" || b.P3 != "0.0" {
		t.Fatalf("percents = %s/%s want 100.0/0.0", b.P2, b.P3)
	}
	requireLabelDistHas(t, b.LabelDist, "unlabeled", 2)

	m := summaryTotalsRE.FindStringSubmatch(rep.Totals)
	if m == nil {
		t.Fatalf("totals line bad: %q", rep.Totals)
	}
	if m[1] != "1" || m[2] != "2" || m[3] != "2" || m[4] != "0" {
		t.Fatalf("totals = %q want roots=1 leaves=2 L2:L3=2:0", rep.Totals)
	}
	if !strings.HasPrefix(rep.Labels, "labels:") {
		t.Fatalf("labels line: %q", rep.Labels)
	}
	requireLabelDistHas(t, strings.TrimPrefix(rep.Labels, "labels:"), "unlabeled", 2)
}
```
