## Expected

- Non-zero exit code.
- Error mentions `-run` and that name-based filters are not supported.
- Suggests path or `--label`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	combined := resp.Stdout + resp.Stderr + resp.ParseErr
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit for -run, got 0\n%s", combined)
	}
	lower := strings.ToLower(combined)
	if !strings.Contains(combined, "-run") {
		t.Fatalf("expected -run in error:\n%s", combined)
	}
	if !strings.Contains(lower, "not supported") {
		t.Fatalf("expected 'not supported' message:\n%s", combined)
	}
	if !strings.Contains(combined, "--label") && !strings.Contains(lower, "path") {
		t.Fatalf("expected path/--label guidance:\n%s", combined)
	}
}
```
