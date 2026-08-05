## Expected

- leaves=2, L2:L3=1:1 (50.0%/50.0%).
- labelDist includes e2e=1 and unlabeled=1.

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
	if b.Leaves != 2 || b.L2 != 1 || b.L3 != 1 {
		t.Fatalf("stats leaves=%d L2:L3=%d:%d want 2 / 1:1", b.Leaves, b.L2, b.L3)
	}
	if b.P2 != "50.0" || b.P3 != "50.0" {
		t.Fatalf("percents %s/%s want 50.0/50.0", b.P2, b.P3)
	}
	requireLabelDistHas(t, b.LabelDist, "e2e", 1)
	requireLabelDistHas(t, b.LabelDist, "unlabeled", 1)
}
```
