## Expected

- The command fails without writing a version.
- The error names `metadata.version`.

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
	if resp.Stdout != "" {
		t.Fatalf("stdout = %q, want empty", resp.Stdout)
	}
	if !strings.Contains(resp.Stderr, "skill tdd has no metadata.version") {
		t.Fatalf("stderr missing version error: %s", resp.Stderr)
	}
}
```
