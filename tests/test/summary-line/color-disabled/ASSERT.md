## Expected

- Command succeeds.
- Summary is plain `PASS (1/1)` with no ANSI escape sequences anywhere in stdout.

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
	summary := findResultSummary(resp.Stdout)
	if strings.TrimSpace(summary) != "PASS (1/1)" {
		t.Fatalf("expected plain PASS (1/1), got %q\nstdout:\n%s", summary, resp.Stdout)
	}
}
```