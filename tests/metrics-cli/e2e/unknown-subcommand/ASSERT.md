---
label: heavy
explanation: L3 product binary smoke for unknown metrics subcommand exit/usage
---

## Expected

- Exit code is non-zero.
- Combined output mentions the unknown command or shows usage/help for metrics.

## Exit Code

- non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit for unknown subcommand; stdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	out := strings.ToLower(combinedOut(resp))
	// Accept common CLI phrasings.
	if !strings.Contains(out, "unknown") &&
		!strings.Contains(out, "not found") &&
		!strings.Contains(out, "invalid") &&
		!strings.Contains(out, "usage") &&
		!strings.Contains(out, "help") {
		t.Fatalf("expected error/usage hint for unknown subcommand:\n%s", combinedOut(resp))
	}
}
```
