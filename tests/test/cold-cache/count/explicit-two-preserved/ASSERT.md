---
label: heavy
---

## Expected

- Command succeeds (exit 0).
- The go test command preview includes `-count=2` (explicit count not forced to 1).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
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
	if !strings.Contains(line, "-count=2") {
		t.Fatalf("expected go test line to preserve -count=2 under --cold-cache, got:\n%s\nfull stderr:\n%s", line, resp.Stderr)
	}
}
```
