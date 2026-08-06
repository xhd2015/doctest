## Expected
- Second `build.Test` run reports `1 Cached` in the summary.
- The `1 Cached` and `0 Fail` segments are wrapped in gray ANSI codes.

## Exit Code
- `err` is nil (cached packages still pass).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.TestErr != nil {
		t.Fatalf("expected cached run to succeed, got: %v\noutput:\n%s", resp.TestErr, resp.Output)
	}
	if !strings.Contains(resp.Summary, "1 Cached") {
		t.Fatalf("expected cached summary, got:\n%s", resp.Output)
	}
	if !metricIsColored(resp.Summary, "1 Cached") {
		t.Fatalf("expected gray 1 Cached segment, summary:\n%s", resp.Summary)
	}
	if !failFieldIsColored(resp.Summary) {
		t.Fatalf("expected gray 0 Fail field when zero, summary:\n%s", resp.Summary)
	}
}
```