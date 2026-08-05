## Expected

- Exit 0.
- stdout contains gray SGR (`\x1b[90m`) on meta (counts, L2:L3, summary, `---`).
- Path field remains plain (no ANSI in path column).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	requireOK(t, resp, err)
	requireGrayMeta(t, resp.Stdout)
	// Still a recognizable report (strip ANSI for structure if needed)
	if !strings.Contains(resp.Stdout, "L2:L3=") {
		t.Fatalf("missing L2:L3 marker:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "---") {
		t.Fatalf("missing summary separator:\n%s", resp.Stdout)
	}
}
```
