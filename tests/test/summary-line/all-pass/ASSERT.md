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
	if !strings.Contains(resp.Stdout, "...  (3 Run, 3 Pass, 0 Fail, 0 Cached)") {
		t.Fatalf("expected inline progress summary, got:\n%s", resp.Stdout)
	}
	if strings.Count(stripANSI(resp.Stdout), ".") < 3 {
		t.Fatalf("expected at least 3 progress dots, got:\n%s", resp.Stdout)
	}
	summary := findResultSummary(resp.Stdout)
	if stripANSI(strings.TrimSpace(summary)) != "PASS (3/3)" {
		t.Fatalf("expected PASS (3/3) summary, got %q\nstdout:\n%s", summary, resp.Stdout)
	}
	summaryIsLastLine(t, resp.Stdout)
	if countResultSummaries(resp.Stdout) != 1 {
		t.Fatalf("expected exactly one result summary line, got %d\nstdout:\n%s",
			countResultSummaries(resp.Stdout), resp.Stdout)
	}
}
```