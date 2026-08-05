## Expected

- Even for a single root, summary footer is present (not body-only).
- stdout ends with a trailing newline after the last summary line.

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
		t.Fatalf("single root must still have summary footer:\n%s", resp.Stdout)
	}
	m := summaryTotalsRE.FindStringSubmatch(rep.Totals)
	if m == nil || m[1] != "1" || m[2] != "1" {
		t.Fatalf("totals %q", rep.Totals)
	}
	if !strings.HasPrefix(rep.Labels, "labels:") {
		t.Fatalf("labels line %q", rep.Labels)
	}
	if !strings.HasSuffix(resp.Stdout, "\n") {
		t.Fatalf("stdout must end with trailing newline")
	}
	// last content line should be labels, then exactly one trailing newline (no extra blank)
	trim := strings.TrimSuffix(resp.Stdout, "\n")
	if strings.HasSuffix(trim, "\n") {
		t.Fatalf("unexpected extra trailing blank line:\n%q", resp.Stdout)
	}
}
```
