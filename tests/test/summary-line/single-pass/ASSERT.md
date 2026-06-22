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
	if !strings.Contains(resp.Stdout, ".  (1 Run, 1 Pass, 0 Fail, 0 Cached)") {
		t.Fatalf("expected inline progress summary, got:\n%s", resp.Stdout)
	}
	summary := findResultSummary(resp.Stdout)
	if stripANSI(strings.TrimSpace(summary)) != "PASS (1/1)" {
		t.Fatalf("expected PASS (1/1) summary, got %q\nstdout:\n%s", summary, resp.Stdout)
	}
	summaryIsLastLine(t, resp.Stdout)
}
```