---
label: heavy
---

## Expected

- Command succeeds (exit 0).
- stdout shows 3 dots and inline `(3 Run, 3 Pass, 0 Fail, 0 Cached)`.
- Final stdout line is `PASS (3/3)` (plain text; pipe disables color).

## Side Effects

- None beyond subprocess output.

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
	if !strings.Contains(plainInline, "(3 Run, 3 Pass, 0 Fail, 0 Cached) in ") {
		t.Fatalf("expected (3 Run, 3 Pass, 0 Fail, 0 Cached) in <duration>, got %q", plainInline)
	}
	inlineIdx := strings.Index(resp.Stdout, "  (")
	if inlineIdx < 0 || strings.Count(resp.Stdout[:inlineIdx], ".") < 3 {
		t.Fatalf("expected at least 3 progress dots before inline summary, got:\n%s", resp.Stdout)
	}
	summary := findResultSummary(resp.Stdout)
	plainSummary := stripANSI(strings.TrimSpace(summary))
	if !strings.HasPrefix(plainSummary, "PASS (3/3) in ") {
		t.Fatalf("expected PASS (3/3) in <duration>, got %q\nstdout:\n%s", plainSummary, resp.Stdout)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v\nsummary: %q", err, summary)
	}
	summaryIsLastLine(t, resp.Stdout)
	if countResultSummaries(resp.Stdout) != 1 {
		t.Fatalf("expected exactly one result summary line, got %d\nstdout:\n%s",
			countResultSummaries(resp.Stdout), resp.Stdout)
	}
}
```