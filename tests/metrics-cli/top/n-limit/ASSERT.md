## Expected

- Exit code 0.
- Top two by elapsed include `group/slow-leaf` and `group/labeled-leaf`.
- `group/fast-leaf` is not listed (truncated by --n 2).

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
	mustContain(t, out, "group/slow-leaf", "top --n")
	// second slowest overall is labeled-leaf (3s)
	mustContain(t, out, "group/labeled-leaf", "top --n")
	if strings.Contains(out, "group/fast-leaf") {
		t.Fatalf("--n 2 should omit fast-leaf:\n%s", out)
	}
}
```
