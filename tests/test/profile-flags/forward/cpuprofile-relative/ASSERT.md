---
label: heavy
explanation: nested go test with cpuprofile; compile + profile overhead
---

## Expected
- Exit code 0.
- stderr shows the go test command with `-cpuprofile=` using an absolute path under WorkDir ending in `profiles/cpu.out`.

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
	wantAbs, err := filepath.Abs(filepath.Join(req.WorkDir, "profiles", "cpu.out"))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	// Prefer equals form like -timeout=; also accept space form.
	if !strings.Contains(resp.Stderr, "-cpuprofile="+wantAbs) &&
		!strings.Contains(resp.Stderr, "-cpuprofile "+wantAbs) {
		t.Fatalf("expected stderr to contain abs -cpuprofile=%s, got:\n%s", wantAbs, resp.Stderr)
	}
	if strings.Contains(resp.Stderr, "-cpuprofile=profiles/cpu.out") ||
		strings.Contains(resp.Stderr, "-cpuprofile profiles/cpu.out") {
		t.Fatalf("relative path must not appear un-resolved on go command line:\n%s", resp.Stderr)
	}
}
```
