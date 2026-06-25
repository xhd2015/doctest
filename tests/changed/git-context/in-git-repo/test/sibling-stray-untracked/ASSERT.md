## Expected

- Exit code 0.
- Summary shows exactly 1 run (leaf_a only).
- Sibling `leaf_b` must be skipped even when it has an unrelated untracked file.

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
	if !strings.Contains(inline, "(1 Run, 1 Pass, 0 Fail") {
		t.Fatalf("expected (1 Run, 1 Pass, 0 Fail...), got %q\nstdout:\n%s", inline, resp.Stdout)
	}
}
```