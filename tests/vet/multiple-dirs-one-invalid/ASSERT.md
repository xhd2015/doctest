## Expected
- The command fails with non-zero exit code.
- stderr contains `root must contain DOCTEST.md` (error from the invalid directory).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected nonzero exit, got 0\nstdout:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stderr, "root must contain DOCTEST.md") {
		t.Fatalf("stderr missing expected error:\n%s", resp.Stderr)
	}
}
```
