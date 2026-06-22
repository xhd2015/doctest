## Expected

- Command succeeds.
- Summary line is `PASS (1/1)` with ANSI green wrapping the entire token.

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
	summary := findResultSummary(resp.Stdout)
	if summary == "" {
		t.Fatalf("expected PASS summary line, got:\n%s", resp.Stdout)
	}
	if stripANSI(strings.TrimSpace(summary)) != "PASS (1/1)" {
		t.Fatalf("expected PASS (1/1), got %q", summary)
	}
	if !containsANSI(summary) {
		t.Fatalf("expected green ANSI on entire PASS (1/1) token, got plain %q", summary)
	}
	if strings.Contains(summary, "FAIL (") {
		t.Fatalf("expected PASS summary, got %q", summary)
	}
}
```