---
label: heavy
---

## Expected
- The command succeeds (exit code 0).
- stdout and stderr contain **no** `t.Logf` output from the generated test (marker absent).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if strings.Contains(resp.Stdout, unwantedNonVerboseLogfMarker) {
		t.Fatalf("stdout must not contain t.Logf output in non-verbose mode, got:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stderr, unwantedNonVerboseLogfMarker) {
		t.Fatalf("stderr must not contain t.Logf output in non-verbose mode, got:\n%s", resp.Stderr)
	}
}
```