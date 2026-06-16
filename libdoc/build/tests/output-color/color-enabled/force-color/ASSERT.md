## Expected
- `build.Test` succeeds.
- Summary contains green-wrapped `1 Pass` and gray-wrapped `0 Fail` segments.
- The single pass dot is plain (no ANSI).

## Exit Code
- `err` is nil.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.TestErr != nil {
		t.Fatalf("expected build.Test to succeed, got: %v", resp.TestErr)
	}
	if !metricIsColored(resp.Summary, "1 Pass") {
		t.Fatalf("expected green 1 Pass segment, summary:\n%s", resp.Summary)
	}
	if !failFieldIsColored(resp.Summary) {
		t.Fatalf("expected gray 0 Fail field when zero, summary:\n%s", resp.Summary)
	}
	if containsANSI(resp.Dots) {
		t.Fatalf("pass dot must stay plain, dots:\n%q", resp.Dots)
	}
	if !strings.Contains(stripANSI(resp.Summary), "1 Run") {
		t.Fatalf("expected 1 Run in summary, got:\n%s", resp.Summary)
	}
}
```