---
label: heavy
explanation: nested go test with -outputdir and -cpuprofile
---

## Expected
- Exit code 0.
- stderr contains abs-resolved `-outputdir=` under WorkDir.
- stderr also contains abs-resolved `-cpuprofile=` (paired flag).

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
	wantOut, err := filepath.Abs(filepath.Join(req.WorkDir, "outprof"))
	if err != nil {
		t.Fatalf("abs outputdir: %v", err)
	}
	if !strings.Contains(resp.Stderr, "-outputdir="+wantOut) &&
		!strings.Contains(resp.Stderr, "-outputdir "+wantOut) {
		t.Fatalf("expected stderr to contain -outputdir=%s, got:\n%s", wantOut, resp.Stderr)
	}
	wantCPU, err := filepath.Abs(filepath.Join(req.WorkDir, "cpu.out"))
	if err != nil {
		t.Fatalf("abs cpuprofile: %v", err)
	}
	if !strings.Contains(resp.Stderr, "-cpuprofile="+wantCPU) &&
		!strings.Contains(resp.Stderr, "-cpuprofile "+wantCPU) {
		t.Fatalf("expected stderr to contain -cpuprofile=%s, got:\n%s", wantCPU, resp.Stderr)
	}
}
```
