## Expected

- Exit 0.
- stdout has no ANSI ESC sequences.

## Exit Code

- 0

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	requireOK(t, resp, err)
	requireNoANSI(t, resp.Stdout, "stdout")
	rep := parseListReport(t, resp.Stdout)
	if len(rep.Body) != 1 {
		t.Fatalf("expected report body:\n%s", resp.Stdout)
	}
}
```
