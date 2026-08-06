## Expected

- Exit code is non-zero (prepare failed for the bad tree).
- stdout final summary is `FAIL (p/t)` — never `PASS (` when prepare failed.
- stderr (CLI error) contains `prepare failed:`.
- Combined output mentions the bad tree path (`bad`).

## Exit Code

- Non-zero.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit when prepare failed, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if strings.Contains(resp.Stdout, "PASS (") {
		t.Fatalf("must not print PASS when prepare failed:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	summary := findResultSummary(resp.Stdout)
	if summary == "" {
		t.Fatalf("expected FAIL (p/t) summary when survivors ran:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	plainSummary := stripANSI(strings.TrimSpace(summary))
	if !strings.HasPrefix(plainSummary, "FAIL (") {
		t.Fatalf("expected FAIL (p/t) summary, got %q\nstdout:\n%s", plainSummary, resp.Stdout)
	}
	combined := resp.Stdout + "\n" + resp.Stderr
	if !strings.Contains(combined, "prepare failed:") {
		t.Fatalf("expected prepare failed: label in process error, got:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(combined, "bad") {
		t.Fatalf("expected error to mention bad tree path, got:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
}
```
