## Expected
- The command succeeds with exit code 0.
- stdout contains `[vet] validating` (directory-level output).
- stdout contains the filename `SETUP.md` (file-level output).

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
		t.Fatalf("expected zero exit, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "[vet] validating") {
		t.Fatalf("stdout missing directory-level verbose output:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "SETUP.md") {
		t.Fatalf("stdout missing file-level verbose output:\n%s", resp.Stdout)
	}
}
```
