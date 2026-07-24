## Expected

- Exit code 0.
- stdout includes `old/only-leaf` and/or stem `showid01`.
- stdout does not need `group/slow-leaf` (should prefer older run).

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
	if !strings.Contains(out, "old/only-leaf") && !strings.Contains(out, "showid01") {
		t.Fatalf("show by id missing older run markers:\n%s", out)
	}
	if strings.Contains(out, "group/slow-leaf") && !strings.Contains(out, "old/only-leaf") {
		t.Fatalf("show by id appears to dump the wrong run:\n%s", out)
	}
}
```
