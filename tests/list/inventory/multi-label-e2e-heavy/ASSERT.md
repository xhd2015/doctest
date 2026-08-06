## Expected

- leaves=1, L2:L3=0:1 (0.0%/100.0%) — e2e identity makes L3.
- labelDist: e2e=1 and slow=1 (multi-label leaf counts once per name).
- unlabeled=0 included in dist (or only labeled names + unlabeled per product; require e2e and slow).

## Exit Code

- 0

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	requireOK(t, resp, err)
	rep := parseListReport(t, resp.Stdout)
	b := findBody(t, rep.Body, req.Roots[0])
	if b.Leaves != 1 || b.L2 != 0 || b.L3 != 1 {
		t.Fatalf("stats leaves=%d L2:L3=%d:%d want 1 / 0:1", b.Leaves, b.L2, b.L3)
	}
	if b.P2 != "0.0" || b.P3 != "100.0" {
		t.Fatalf("percents %s/%s want 0.0/100.0", b.P2, b.P3)
	}
	requireLabelDistHas(t, b.LabelDist, "e2e", 1)
	requireLabelDistHas(t, b.LabelDist, "slow", 1)
	// Multi-label leaf is not unlabeled
	requireLabelDistHas(t, b.LabelDist, "unlabeled", 0)
}
```
