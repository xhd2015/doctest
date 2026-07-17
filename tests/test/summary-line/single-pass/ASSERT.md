---
label: heavy
---

## Expected

- Command succeeds.
- stdout shows 1 dot, inline `(1 Run, 1 Pass, 0 Fail, 0 Cached)`, and final line `PASS (1/1)`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	inline := findInlineSummaryLine(resp.Stdout)
	if inline == "" {
		t.Fatalf("expected inline progress summary with duration, got:\n%s", resp.Stdout)
	}
	plainInline := stripANSI(strings.TrimSpace(inline))
	if !strings.Contains(plainInline, "(1 Run, 1 Pass, 0 Fail, 0 Cached) in ") {
		t.Fatalf("expected (1 Run, 1 Pass, 0 Fail, 0 Cached) in <duration>, got %q", plainInline)
	}
	summary := findResultSummary(resp.Stdout)
	plainSummary := stripANSI(strings.TrimSpace(summary))
	if !strings.HasPrefix(plainSummary, "PASS (1/1) in ") {
		t.Fatalf("expected PASS (1/1) in <duration>, got %q\nstdout:\n%s", plainSummary, resp.Stdout)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v\nsummary: %q", err, summary)
	}
	summaryIsLastLine(t, resp.Stdout)
}
```