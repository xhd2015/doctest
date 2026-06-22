## Expected

- Exit code non-zero.
- Summary line is `FAIL (0/1)` with ANSI red wrapping the entire token.

## Exit Code

- Non-zero.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	summary := findResultSummary(resp.Stdout)
	if summary == "" {
		t.Fatalf("expected FAIL summary line, got:\n%s", resp.Stdout)
	}
	if stripANSI(strings.TrimSpace(summary)) != "FAIL (0/1)" {
		t.Fatalf("expected FAIL (0/1), got %q", summary)
	}
	if !containsANSI(summary) {
		t.Fatalf("expected red ANSI on entire FAIL (0/1) token, got plain %q", summary)
	}
}
```