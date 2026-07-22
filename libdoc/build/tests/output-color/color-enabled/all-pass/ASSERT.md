---
label: heavy
---

## Expected
- `build.Test` succeeds with two plain dots before the summary.
- `2 Pass` is colored; `0 Fail` and `0 Cached` are gray; `2 Run` is plain.

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

	plainDots := stripANSI(resp.Dots)
	if strings.Count(plainDots, ".") != 2 {
		t.Fatalf("expected 2 plain dots, got %q (raw: %q)", plainDots, resp.Dots)
	}
	if containsANSI(resp.Dots) {
		t.Fatalf("pass dots must not be colored, dots:\n%q", resp.Dots)
	}

	if !metricIsColored(resp.Summary, "2 Pass") {
		t.Fatalf("expected green 2 Pass, summary:\n%s", resp.Summary)
	}
	if !failFieldIsColored(resp.Summary) {
		t.Fatalf("expected gray 0 Fail field when zero, summary:\n%s", resp.Summary)
	}
	if !metricIsColored(resp.Summary, "0 Cached") {
		t.Fatalf("expected gray 0 Cached even when zero, summary:\n%s", resp.Summary)
	}
	if !metricIsPlain(resp.Summary, "2 Run") {
		t.Fatalf("expected plain 2 Run, summary:\n%s", resp.Summary)
	}
}
```