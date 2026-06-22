## Expected
- Command exits 0 (no tests is printed to stderr but not an error exit).
- Stderr contains "no tests".
- No test cases are run.

## Exit Code
- Exit code 0.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr:\n%s\nstdout:\n%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}

	combined := resp.Stderr + "\n" + resp.Stdout
	if !strings.Contains(combined, "no tests") {
		t.Fatalf("expected 'no tests' in output:\nstderr:\n%s\nstdout:\n%s", resp.Stderr, resp.Stdout)
	}
}
```
