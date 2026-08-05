## Expected

- Exit 0.
- Exactly two body lines (deduped).
- Paths sorted.
- Summary roots=2.

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
		t.Fatalf("want 2 deduped roots, got %d\n%s", len(rep.Body), resp.Stdout)
	}
	requireSortedPaths(t, rep.Body)
	findBody(t, rep.Body, req.Roots[0])
	findBody(t, rep.Body, req.Roots[1])
	m := summaryTotalsRE.FindStringSubmatch(rep.Totals)
	if m == nil || m[1] != "2" {
		t.Fatalf("totals = %q want roots=2", rep.Totals)
	}
}
```
