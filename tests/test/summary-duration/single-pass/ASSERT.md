---
label: heavy
---

## Expected

- Command succeeds.
- stdout shows 1 dot, inline `(1 Run, 1 Pass, 0 Fail, 0 Cached) in <duration>`, and final `PASS (1/1) in <duration>`.
- Both durations parse as valid Go duration strings.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
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
	if !inlineSummaryPlainRe.MatchString(plainInline) {
		t.Fatalf("inline summary must match (N Run, N Pass, N Fail, N Cached) in DURATION, got %q", plainInline)
	}
	if _, err := parseInlineSummaryDuration(inline); err != nil {
		t.Fatalf("inline duration must parse: %v\nline: %q", err, inline)
	}
	summary := findResultSummary(resp.Stdout)
	if summary == "" {
		t.Fatalf("expected final PASS summary with duration, got:\n%s", resp.Stdout)
	}
	plainSummary := stripANSI(strings.TrimSpace(summary))
	if !finalSummaryPlainRe.MatchString(plainSummary) {
		t.Fatalf("expected PASS (1/1) in <duration>, got %q", plainSummary)
	}
	if !strings.HasPrefix(plainSummary, "PASS (1/1) in ") {
		t.Fatalf("expected PASS (1/1) in <duration>, got %q", plainSummary)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v\nsummary: %q", err, summary)
	}
	summaryIsLastLine(t, resp.Stdout)
}
```