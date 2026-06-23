## Expected

- Exit code 0.
- Root `DOCTEST.md` version check is skipped because the file is unchanged.

## Exit Code

0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if strings.Contains(resp.Stderr, "## Version") {
		t.Fatalf("stderr should not report missing Version (root was skipped):\n%s", resp.Stderr)
	}
}
```