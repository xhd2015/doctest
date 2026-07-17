---
label: heavy
---

## Expected

- Exit code is non-zero (at least one failure).
- stdout contains `FAIL\t` lines before the summary.
- Final stdout line is `FAIL (2/3)`.
- stdout includes inline `(3 Run, 2 Pass, 1 Fail, 0 Cached)` after dots.

## Exit Code

- Non-zero.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	inline := findInlineSummaryLine(resp.Stdout)
	if inline == "" {
		t.Fatalf("expected inline progress summary with duration, got:\n%s", resp.Stdout)
	}
	plainInline := stripANSI(strings.TrimSpace(inline))
	if !strings.Contains(plainInline, "(3 Run, 2 Pass, 1 Fail, 0 Cached) in ") {
		t.Fatalf("expected (3 Run, 2 Pass, 1 Fail, 0 Cached) in <duration>, got %q", plainInline)
	}

	hasFailTab := false
	summaryIdx := -1
	for i, line := range strings.Split(resp.Stdout, "\n") {
		if strings.HasPrefix(line, "FAIL\t") {
			hasFailTab = true
		}
		plain := strings.TrimSpace(stripANSI(line))
		if strings.HasPrefix(plain, "FAIL (") || strings.HasPrefix(plain, "PASS (") {
			summaryIdx = i
		}
	}
	if !hasFailTab {
		t.Fatalf("expected FAIL\\t lines before summary, got:\n%s", resp.Stdout)
	}

	summary := findResultSummary(resp.Stdout)
	plainSummary := stripANSI(strings.TrimSpace(summary))
	if !strings.HasPrefix(plainSummary, "FAIL (2/3) in ") {
		t.Fatalf("expected FAIL (2/3) in <duration>, got %q\nstdout:\n%s", plainSummary, resp.Stdout)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v\nsummary: %q", err, summary)
	}
	summaryIsLastLine(t, resp.Stdout)

	if summaryIdx >= 0 {
		before := strings.Join(strings.Split(resp.Stdout, "\n")[:summaryIdx], "\n")
		if !strings.Contains(before, "FAIL\t") {
			t.Fatalf("FAIL\\t lines must appear before summary line\nstdout:\n%s", resp.Stdout)
		}
	}
}
```