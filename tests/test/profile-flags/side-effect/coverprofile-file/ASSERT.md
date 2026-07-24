---
label: heavy
explanation: nested go test writes coverprofile file
---

## Expected
- Exit code 0.
- The coverprofile path from Args exists on disk after the run.

## Side Effects
- Coverage profile file created at the absolute `-coverprofile` path.

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
	var covPath string
	for i, a := range req.Args {
		if a == "-coverprofile" && i+1 < len(req.Args) {
			covPath = req.Args[i+1]
			break
		}
	}
	if covPath == "" {
		t.Fatal("missing -coverprofile path in req.Args")
	}
	info, statErr := os.Stat(covPath)
	if statErr != nil {
		t.Fatalf("expected coverprofile file at %s: %v\nstdout:\n%s\nstderr:\n%s", covPath, statErr, resp.Stdout, resp.Stderr)
	}
	if info.Size() == 0 {
		t.Fatalf("coverprofile file %s is empty", covPath)
	}
}
```
