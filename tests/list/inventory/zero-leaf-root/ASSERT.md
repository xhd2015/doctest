## Expected

- One body line: leaves=0, L2:L3=0:0, **no** percent group.
- Summary roots=1 leaves=0 L2:L3=0:0 without percent group.
- labelDist includes unlabeled=0 (or empty-only unlabeled).

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
	b := findBody(t, rep.Body, req.Roots[0])
	if b.Leaves != 0 || b.L2 != 0 || b.L3 != 0 {
		t.Fatalf("stats leaves=%d L2:L3=%d:%d want 0 / 0:0", b.Leaves, b.L2, b.L3)
	}
	if b.HasPct {
		t.Fatalf("percent group must be omitted when leaves==0: %+v", b)
	}
	// Body line must not contain "(%.%/" style percent group
	for _, ln := range strings.Split(resp.Stdout, "\n") {
		if strings.Contains(ln, b.Path) && strings.Contains(ln, "%") {
			t.Fatalf("zero-leaf body must not include percent tokens: %q", ln)
		}
	}
	m := summaryTotalsRE.FindStringSubmatch(rep.Totals)
	if m == nil || m[1] != "1" || m[2] != "0" || m[3] != "0" || m[4] != "0" {
		t.Fatalf("totals %q want roots=1 leaves=0 L2:L3=0:0", rep.Totals)
	}
	if m[5] != "" || m[6] != "" {
		t.Fatalf("summary must omit percent group when leaves==0: %q", rep.Totals)
	}
	requireLabelDistHas(t, b.LabelDist, "unlabeled", 0)
}
```
