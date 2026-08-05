## Expected

- Exit 0.
- Both parent and nested child appear as roots.
- Parent leaf count is 1 (only `own`), not 2.
- Child leaf count is 1.
- Summary leaves = 2 (sum of ownership).

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
	if len(rep.Body) != 2 {
		t.Fatalf("want 2 roots, got %d\n%s", len(rep.Body), resp.Stdout)
	}
	parent := findBody(t, rep.Body, req.Roots[0])
	child := findBody(t, rep.Body, req.Roots[1])
	if parent.Leaves != 1 {
		t.Fatalf("parent must own 1 leaf (exclude nested), got %d\n%s", parent.Leaves, resp.Stdout)
	}
	if child.Leaves != 1 {
		t.Fatalf("child must own 1 leaf, got %d", child.Leaves)
	}
	m := summaryTotalsRE.FindStringSubmatch(rep.Totals)
	if m == nil || m[1] != "2" || m[2] != "2" {
		t.Fatalf("totals = %q want roots=2 leaves=2", rep.Totals)
	}
}
```
