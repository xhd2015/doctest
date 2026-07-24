---
label: heavy
---

## Expected

- Command succeeds.
- Inline and final parsed durations are each at least 1 second.

```go
import (
	"strings"
	"testing"
	"time"
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
		t.Fatalf("expected inline summary with duration, got:\n%s", resp.Stdout)
	}
	inlineDur, err := parseInlineSummaryDuration(inline)
	if err != nil {
		t.Fatalf("inline duration must parse: %v", err)
	}
	if inlineDur < time.Second {
		t.Fatalf("expected inline duration >= 1s, got %v\nline: %q", inlineDur, stripANSI(inline))
	}
	finalDur, err := parseFinalSummaryDuration(resp.Stdout)
	if err != nil {
		t.Fatalf("final duration must parse: %v", err)
	}
	if finalDur < time.Second {
		t.Fatalf("expected final duration >= 1s, got %v", finalDur)
	}
	summary := findResultSummary(resp.Stdout)
	plainSummary := stripANSI(strings.TrimSpace(summary))
	if !strings.HasPrefix(plainSummary, "PASS (1/1) in ") {
		t.Fatalf("expected PASS (1/1) in <duration>, got %q", plainSummary)
	}
	summaryIsLastLine(t, resp.Stdout)
}
```