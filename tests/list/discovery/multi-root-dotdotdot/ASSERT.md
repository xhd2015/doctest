## Expected

- Exit 0.
- Two body lines (alpha and beta), paths sorted.
- Each root reports leaves=1.
- Summary roots=2 leaves=2.
- Body lines appear before the `---` summary block.

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
	if len(rep.Body) != 2 {
		t.Fatalf("want 2 body lines, got %d\n%s", len(rep.Body), resp.Stdout)
	}
	requireSortedPaths(t, rep.Body)
	if !rep.HasSep {
		t.Fatalf("expected summary separator:\n%s", resp.Stdout)
	}
	for _, root := range req.Roots {
		b := findBody(t, rep.Body, root)
		if b.Leaves != 1 || b.L2 != 1 || b.L3 != 0 {
			t.Fatalf("%s: leaves/L2/L3 = %d/%d/%d want 1/1/0", root, b.Leaves, b.L2, b.L3)
		}
	}
	m := summaryTotalsRE.FindStringSubmatch(rep.Totals)
	if m == nil || m[1] != "2" || m[2] != "2" {
		t.Fatalf("totals = %q want roots=2 leaves=2", rep.Totals)
	}
	// Streaming contract (observable): first non-empty content is a body path, not ---.
	trim := strings.TrimPrefix(resp.Stdout, "\n")
	if strings.HasPrefix(strings.TrimSpace(trim), "---") {
		t.Fatalf("body must precede summary; stdout starts with ---:\n%s", resp.Stdout)
	}
}
```
