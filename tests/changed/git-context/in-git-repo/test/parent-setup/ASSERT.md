---
label: heavy
---

## Expected

- Exit code 0.
- Summary shows exactly 2 runs (both leaves under shared parent).

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
	// Unified suite: one go test package, two leaf subtests → summary is
	// leaf-level (2 Pass) or package-level (1 Pass) depending on reporter;
	// always require both leaves in PASS (2/2) and no Fail.
	if strings.Contains(inline, "Fail") && !strings.Contains(inline, "0 Fail") {
		t.Fatalf("unexpected failures in summary %q\nstdout:\n%s", inline, resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "PASS (2/2)") && !strings.Contains(inline, "2 Pass") {
		t.Fatalf("expected both leaves to pass (PASS (2/2) or 2 Pass), got %q\nstdout:\n%s", inline, resp.Stdout)
	}
}
```