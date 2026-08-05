## Expected

- Two body lines.
- After bodies: blank line, `---`, totals, labels.
- Summary: roots=2 leaves=3 L2:L3=2:1 (L2 66.7% / L3 33.3%).
- labels: e2e=1 unlabeled=2 (order: count desc then name).
- Body L2/L3/leaves sum to summary.

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
	if !rep.HasSep {
		t.Fatalf("expected blank+--- summary:\n%s", resp.Stdout)
	}
	if len(rep.Body) != 2 {
		t.Fatalf("want 2 body lines, got %d\n%s", len(rep.Body), resp.Stdout)
	}
	sumLeaves, sumL2, sumL3 := 0, 0, 0
	for _, b := range rep.Body {
		sumLeaves += b.Leaves
		sumL2 += b.L2
		sumL3 += b.L3
	}
	if sumLeaves != 3 || sumL2 != 2 || sumL3 != 1 {
		t.Fatalf("body sums leaves=%d L2=%d L3=%d want 3/2/1", sumLeaves, sumL2, sumL3)
	}
	m := summaryTotalsRE.FindStringSubmatch(rep.Totals)
	if m == nil {
		t.Fatalf("totals line: %q", rep.Totals)
	}
	if m[1] != "2" || m[2] != "3" || m[3] != "2" || m[4] != "1" {
		t.Fatalf("totals %q want roots=2 leaves=3 L2:L3=2:1", rep.Totals)
	}
	if m[5] != "66.7" || m[6] != "33.3" {
		t.Fatalf("summary percents L2=%s L3=%s want 66.7 / 33.3", m[5], m[6])
	}
	lab := strings.TrimPrefix(rep.Labels, "labels:")
	requireLabelDistHas(t, lab, "e2e", 1)
	requireLabelDistHas(t, lab, "unlabeled", 2)
	// Structure: blank line immediately before ---
	if !strings.Contains(resp.Stdout, "\n\n---\n") {
		t.Fatalf("expected blank line before ---:\n%q", resp.Stdout)
	}
}
```
