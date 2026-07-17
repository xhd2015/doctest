---
label: heavy
explanation: nested go test with memprofilerate=0
---

## Expected
- Exit code 0.
- stderr contains `-memprofilerate=0` or `-memprofilerate 0` (exact zero, not omitted).

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
		t.Fatalf("expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	hasRate := strings.Contains(resp.Stderr, "-memprofilerate=0") ||
		strings.Contains(resp.Stderr, "-memprofilerate 0")
	if !hasRate {
		t.Fatalf("expected stderr to contain -memprofilerate=0, got:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "-memprofile=") && !strings.Contains(resp.Stderr, "-memprofile ") {
		t.Fatalf("expected -memprofile on go command line, got:\n%s", resp.Stderr)
	}
}
```
