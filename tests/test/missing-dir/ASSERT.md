## Expected

- The command fails.
- stderr reports that test requires a directory.
- Nonzero exit code.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected nonzero exit, stdout:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stderr, "test requires <dir>") {
		t.Fatalf("stderr missing required-dir message:\n%s", resp.Stderr)
	}
}
```
