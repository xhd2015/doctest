## Expected
- `build.Test` returns an error (one failing leaf).
- Dots and summary contain no ANSI escape sequences.

## Errors
- `resp.TestErr` is non-nil (go test failed).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if resp.TestErr == nil {
		t.Fatal("expected build.Test to fail for failing leaf")
	}
	if !strings.Contains(resp.Summary, "1 Fail") {
		t.Fatalf("expected summary with 1 Fail, got:\n%s", resp.Output)
	}
	if containsANSI(resp.Dots) {
		t.Fatalf("ColorNever must not color fail dot, dots:\n%q", resp.Dots)
	}
	if containsANSI(resp.Summary) {
		t.Fatalf("ColorNever must not color summary, summary:\n%q", resp.Summary)
	}
}
```