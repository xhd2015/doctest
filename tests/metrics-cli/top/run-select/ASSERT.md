## Expected

- Exit code 0.
- stdout includes `old/only-leaf`.
- stdout does not include `group/slow-leaf` (that path is only on the newer run).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	out := resp.Stdout
	mustContain(t, out, "old/only-leaf", "top --run")
	if strings.Contains(out, "group/slow-leaf") {
		t.Fatalf("--run older should not show newer run leaves:\n%s", out)
	}
}
```
