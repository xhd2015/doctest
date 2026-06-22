## Expected

- Command succeeds.
- Summary is plain `PASS (1/1) in <duration>` with no ANSI escape sequences anywhere in stdout.
- Inline summary duration is also plain text.

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
	if containsANSI(resp.Stdout) {
		t.Fatalf("stdout must not contain ANSI with --no-color, got:\n%s", resp.Stdout)
	}
	inline := findInlineSummaryLine(resp.Stdout)
	if inline == "" {
		t.Fatalf("expected inline summary with duration, got:\n%s", resp.Stdout)
	}
	if _, err := parseInlineSummaryDuration(inline); err != nil {
		t.Fatalf("inline duration must parse: %v", err)
	}
	summary := findResultSummary(resp.Stdout)
	plainSummary := strings.TrimSpace(summary)
	if !strings.HasPrefix(plainSummary, "PASS (1/1) in ") {
		t.Fatalf("expected plain PASS (1/1) in <duration>, got %q\nstdout:\n%s", plainSummary, resp.Stdout)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v", err)
	}
}
```