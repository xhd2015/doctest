## Expected

- Exit code non-zero.
- Combined output indicates no runs / not found / empty metrics (case-insensitive keywords).

## Exit Code

- non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero when no runs; stdout:\n%s", resp.Stdout)
	}
	out := strings.ToLower(combinedOut(resp))
	ok := strings.Contains(out, "no run") ||
		strings.Contains(out, "no metrics") ||
		strings.Contains(out, "not found") ||
		strings.Contains(out, "empty") ||
		strings.Contains(out, "no such") ||
		(strings.Contains(out, "run") && strings.Contains(out, "found"))
	if !ok {
		t.Fatalf("expected no-runs messaging:\n%s", combinedOut(resp))
	}
}
```
