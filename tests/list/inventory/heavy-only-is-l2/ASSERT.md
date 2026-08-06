## Expected

- leaves=1, L2:L3=1:0 (100.0%/0.0%) — cost label alone is not L3.
- labelDist includes slow=1.

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
	if b.Leaves != 1 || b.L2 != 1 || b.L3 != 0 {
		t.Fatalf("stats leaves=%d L2:L3=%d:%d want 1 / 1:0 (slow is L2)", b.Leaves, b.L2, b.L3)
	}
	if b.P2 != "100.0" || b.P3 != "0.0" {
		t.Fatalf("percents %s/%s", b.P2, b.P3)
	}
	requireLabelDistHas(t, b.LabelDist, "slow", 1)
	requireLabelDistHas(t, b.LabelDist, "unlabeled", 0)
}
```
