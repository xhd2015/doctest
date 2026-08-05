## Expected

- Exit 0.
- stdout has no ANSI ESC (`\x1b`) — buffers are non-TTY so auto color is off.

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
	// Successful report still has structure
	rep := parseListReport(t, resp.Stdout)
	if len(rep.Body) != 1 || !rep.HasSep {
		t.Fatalf("expected body+summary without ANSI:\n%s", resp.Stdout)
	}
}
```
