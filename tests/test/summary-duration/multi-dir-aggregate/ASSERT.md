---
label: heavy
---

## Expected

- Command succeeds.
- Per-dir inline summaries with duration appear for both directories.
- Exactly one `PASS (3/3) in <duration>` summary at the very end of stdout.

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
	if countInlineSummaries(resp.Stdout) < 2 {
		t.Fatalf("expected per-dir inline summaries with duration for both dirs, got %d\nstdout:\n%s",
			countInlineSummaries(resp.Stdout), resp.Stdout)
	}
	if countResultSummaries(resp.Stdout) != 1 {
		t.Fatalf("expected exactly one aggregated summary, got %d\nstdout:\n%s",
			countResultSummaries(resp.Stdout), resp.Stdout)
	}
	summary := findResultSummary(resp.Stdout)
	plainSummary := stripANSI(strings.TrimSpace(summary))
	if !strings.HasPrefix(plainSummary, "PASS (3/3) in ") {
		t.Fatalf("expected PASS (3/3) in <duration> aggregated summary, got %q\nstdout:\n%s",
			plainSummary, resp.Stdout)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v\nsummary: %q", err, summary)
	}
	summaryIsLastLine(t, resp.Stdout)
}
```