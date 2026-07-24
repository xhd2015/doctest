---
label: heavy
---

## Expected

- Command succeeds.
- Summary line is `PASS (1/1)` with ANSI green wrapping the entire token.

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
	summary := findResultSummary(resp.Stdout)
	if summary == "" {
		t.Fatalf("expected PASS summary line, got:\n%s", resp.Stdout)
	}
	plainSummary := stripANSI(strings.TrimSpace(summary))
	if !strings.HasPrefix(plainSummary, "PASS (1/1) in ") {
		t.Fatalf("expected PASS (1/1) in <duration>, got %q", plainSummary)
	}
	if !finalSummaryPassTokenIsColored(summary) {
		t.Fatalf("expected green ANSI on PASS (1/1) token, got %q", summary)
	}
	if !finalSummaryDurationIsPlain(summary) {
		t.Fatalf("expected plain duration suffix after ' in ', got %q", summary)
	}
	if strings.Contains(summary, "FAIL (") {
		t.Fatalf("expected PASS summary, got %q", summary)
	}
}
```