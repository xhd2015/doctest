## Expected

- Exit code 0.
- stderr contains `no tests` (not `no tests found`).
- stdout has no `PASS (` / `FAIL (` summary line.

## Exit Code

- 0.

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
	if !strings.Contains(resp.Stderr, "no tests") {
		t.Fatalf("stderr must contain 'no tests', got:\n%s", resp.Stderr)
	}
	if strings.Contains(resp.Stderr, "no tests found") {
		t.Fatalf("stderr must not contain old 'no tests found' message:\n%s", resp.Stderr)
	}
	if findResultSummary(resp.Stdout) != "" {
		t.Fatalf("stdout must not contain PASS/FAIL summary when no cases, got:\n%s", resp.Stdout)
	}
}
```