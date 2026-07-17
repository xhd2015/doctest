---
label: heavy
---

## Expected

- Command succeeds.
- stderr shows the go test command includes `-v`.
- stdout has no inline `(N Run, ...)` summary.
- stdout ends with `PASS (1/1) in <duration>` after verbose go-test output.

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
	if !strings.Contains(resp.Stderr, "go test -v") {
		t.Fatalf("expected stderr to contain 'go test -v', got:\n%s", resp.Stderr)
	}
	if findInlineSummaryLine(resp.Stdout) != "" {
		t.Fatalf("verbose mode must not print inline dot summary, got:\n%s", resp.Stdout)
	}
	summary := findResultSummary(resp.Stdout)
	plainSummary := stripANSI(strings.TrimSpace(summary))
	if !strings.HasPrefix(plainSummary, "PASS (1/1) in ") {
		t.Fatalf("expected PASS (1/1) in <duration> after verbose output, got %q\nstdout:\n%s",
			plainSummary, resp.Stdout)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v\nsummary: %q", err, summary)
	}
	summaryIsLastLine(t, resp.Stdout)
}
```