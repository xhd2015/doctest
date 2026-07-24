---
label: heavy
explanation: nested go test writes cpuprofile file
---

## Expected
- Exit code 0.
- The cpuprofile path from Args exists on disk after the run.

## Side Effects
- CPU profile file created at the absolute `-cpuprofile` path.

```go
import (
	"os"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	var cpuPath string
	for i, a := range req.Args {
		if a == "-cpuprofile" && i+1 < len(req.Args) {
			cpuPath = req.Args[i+1]
			break
		}
	}
	if cpuPath == "" {
		t.Fatal("missing -cpuprofile path in req.Args")
	}
	info, statErr := os.Stat(cpuPath)
	if statErr != nil {
		t.Fatalf("expected cpuprofile file at %s: %v\nstdout:\n%s\nstderr:\n%s", cpuPath, statErr, resp.Stdout, resp.Stderr)
	}
	if info.Size() == 0 {
		// Prefer non-empty; go usually writes a header at minimum.
		t.Fatalf("cpuprofile file %s is empty", cpuPath)
	}
}
```
