## Expected

- Exit code 0.
- stdout includes `group/slow-leaf` from the default-suite run.
- stdout does **not** list `only/labeled-suite-leaf` as a top row (that leaf only exists on the non-default suite run).

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
	mustContain(t, out, "group/slow-leaf", "default-only top")
	if strings.Contains(out, "only/labeled-suite-leaf") {
		t.Fatalf("--default-only should not rank labeled-suite-only leaf:\n%s", out)
	}
}
```
