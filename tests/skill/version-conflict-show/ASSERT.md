## Expected

- The command fails before loading skill content.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("harness failed: %v", err)
	}
	if resp.ExitCode != 1 {
		t.Fatalf("exit code = %d, want 1", resp.ExitCode)
	}
	if !strings.Contains(resp.Stderr, "expected exactly one of --show, --install, --list, or --version") {
		t.Fatalf("stderr missing action conflict: %s", resp.Stderr)
	}
}
```
