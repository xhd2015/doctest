## Expected

- Exit code 0.
- Summary shows exactly 2 runs (all leaves affected by root change).

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
	inline := findInlineSummaryLine(resp.Stdout)
	if inline == "" {
		t.Fatalf("missing summary line in stdout:\n%s", resp.Stdout)
	}
	if !strings.Contains(inline, "(2 Run, 2 Pass, 0 Fail") {
		t.Fatalf("expected (2 Run, 2 Pass, 0 Fail...), got %q\nstdout:\n%s", inline, resp.Stdout)
	}
}
```