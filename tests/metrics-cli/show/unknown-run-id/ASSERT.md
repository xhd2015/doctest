## Expected

- Exit code non-zero.
- Combined output mentions not found / unknown / no such run (flexible wording).

## Exit Code

- non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero for unknown run id; stdout:\n%s", resp.Stdout)
	}
	out := strings.ToLower(combinedOut(resp))
	if !strings.Contains(out, "not found") &&
		!strings.Contains(out, "unknown") &&
		!strings.Contains(out, "no such") &&
		!strings.Contains(out, "missing") &&
		!strings.Contains(out, "no run") {
		t.Fatalf("expected not-found messaging:\n%s", combinedOut(resp))
	}
}
```
