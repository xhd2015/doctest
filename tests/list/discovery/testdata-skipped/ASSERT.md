## Expected

- Exit 0.
- Exactly one body line (the real root), not the testdata decoy.
- leaves=1 (real leaf only).

## Exit Code

- 0

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	requireOK(t, resp, err)
	rep := parseListReport(t, resp.Stdout)
	if len(rep.Body) != 1 {
		t.Fatalf("want 1 root (testdata skipped), got %d\n%s", len(rep.Body), resp.Stdout)
	}
	b := findBody(t, rep.Body, req.Roots[0])
	if b.Leaves != 1 {
		t.Fatalf("leaves=%d want 1 (testdata ASSERT ignored)", b.Leaves)
	}
	// No path component .../testdata/... as a listed root
	for _, row := range rep.Body {
		slash := filepath.ToSlash(row.Path)
		if strings.Contains(slash, "/testdata/") || strings.HasSuffix(slash, "/testdata") {
			t.Fatalf("testdata must not appear as root: %q", row.Path)
		}
	}
}
```
