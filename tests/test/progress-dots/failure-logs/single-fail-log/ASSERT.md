## Expected
- stdout contains `SINGLE_FAIL_LOG_MARKER` after the summary (detailed go test failure output forwarded).
- Summary shows one failure: `(1 Run, 0 Pass, 1 Fail, 0 Cached)`.
- Exit code is non-zero.

## Exit Code
- `resp.ExitCode != 0`.

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
		t.Fatalf("expected non-zero exit for failing test, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, singleFailLogMarker) {
		t.Fatalf("stdout must contain failure marker %q\nstdout:\n%s\nstderr:\n%s",
			singleFailLogMarker, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, ".  (1 Run, 0 Pass, 1 Fail, 0 Cached)") {
		t.Fatalf("expected summary line in stdout, got:\n%s", resp.Stdout)
	}
}
```