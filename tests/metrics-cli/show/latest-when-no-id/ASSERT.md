## Expected

- Exit code 0.
- stdout includes a leaf path from the newer run (`group/slow-leaf`) and/or stem `shownew1`.
- Prefer not exclusive `old/only-leaf` without newer markers.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	out := resp.Stdout
	if !strings.Contains(out, "group/slow-leaf") && !strings.Contains(out, "shownew1") {
		t.Fatalf("show latest missing newer run markers:\n%s", out)
	}
}
```
