---
label: heavy
explanation: nested go test with -trace
---

## Expected
- Exit code 0.
- stderr contains abs-resolved `-trace=` under WorkDir.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	wantAbs, err := filepath.Abs(filepath.Join(req.WorkDir, "traces", "run.out"))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	if !strings.Contains(resp.Stderr, "-trace="+wantAbs) &&
		!strings.Contains(resp.Stderr, "-trace "+wantAbs) {
		t.Fatalf("expected stderr to contain -trace=%s, got:\n%s", wantAbs, resp.Stderr)
	}
}
```
