## Expected

- Exit code 0.
- stdout includes unlabeled `group/slow-leaf`.
- stdout does **not** include `group/labeled-leaf`.

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
	mustContain(t, out, "group/slow-leaf", "unlabeled-only")
	if strings.Contains(out, "group/labeled-leaf") {
		t.Fatalf("--unlabeled-only should omit labeled-leaf:\n%s", out)
	}
}
```
