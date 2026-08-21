## Expected

- The command fails because `--header` belongs only to `--show`.

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
	if !strings.Contains(resp.Stderr, "--header is only valid with --show") {
		t.Fatalf("stderr missing header constraint: %s", resp.Stderr)
	}
}
```
