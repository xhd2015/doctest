## Expected

- Command succeeds.
- Per-dir headers and dots appear for both directories.
- Exactly one `PASS (3/3)` summary at the very end of stdout.

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
	if countResultSummaries(resp.Stdout) != 1 {
		t.Fatalf("expected exactly one aggregated summary, got %d\nstdout:\n%s",
			countResultSummaries(resp.Stdout), resp.Stdout)
	}
	summary := findResultSummary(resp.Stdout)
	if stripANSI(strings.TrimSpace(summary)) != "PASS (3/3)" {
		t.Fatalf("expected PASS (3/3) aggregated summary, got %q\nstdout:\n%s",
			summary, resp.Stdout)
	}
	summaryIsLastLine(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, " Run, ") {
		t.Fatalf("expected per-dir inline progress summaries, got:\n%s", resp.Stdout)
	}
}
```