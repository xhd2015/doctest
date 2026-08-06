---
label: e2e
---

## Expected

- Command succeeds (exit 0).
- The go test command preview on stderr includes `-count=1` (forced because count was unset).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	line := goTestCmdLine(resp.Stderr)
	if line == "" {
		t.Fatalf("expected 'cd … && go test …' preview on stderr:\n%s", resp.Stderr)
	}
	// Require -count=1 as a flag token (avoid matching -count=10 etc. via prefix alone).
	if !strings.Contains(line, "-count=1") {
		t.Fatalf("expected go test line to include -count=1 when count unset under --cold-cache, got:\n%s\nfull stderr:\n%s", line, resp.Stderr)
	}
	// Must not only force a different count.
	if strings.Contains(line, "-count=2") {
		t.Fatalf("unexpected -count=2 when count was unset:\n%s", line)
	}
}
```
