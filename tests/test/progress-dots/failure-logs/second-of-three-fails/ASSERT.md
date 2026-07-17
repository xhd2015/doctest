---
label: heavy
---

## Expected
- stdout contains `SECOND_FAIL_LOG_MARKER`.
- stdout has exactly one `--- FAIL:` block (no extra failure output from passing packages).
- stdout has exactly one `FAIL\t` summary line from go test.
- Summary shows `(3 Run, 2 Pass, 1 Fail, 0 Cached)`.
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
		t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, secondFailLogMarker) {
		t.Fatalf("stdout must contain failure marker %q\nstdout:\n%s\nstderr:\n%s",
			secondFailLogMarker, resp.Stdout, resp.Stderr)
	}
	if strings.Count(resp.Stdout, "--- FAIL:") != 1 {
		t.Fatalf("expected exactly one --- FAIL: block, got %d\nstdout:\n%s",
			strings.Count(resp.Stdout, "--- FAIL:"), resp.Stdout)
	}
	failTabCount := 0
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.HasPrefix(line, "FAIL\t") {
			failTabCount++
		}
		if strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "ok\t") {
			t.Fatalf("stdout must not contain raw go test ok lines, got: %q\nfull stdout:\n%s",
				line, resp.Stdout)
		}
	}
	if failTabCount != 1 {
		t.Fatalf("expected exactly one FAIL\\t line, got %d\nstdout:\n%s", failTabCount, resp.Stdout)
	}
	inline := findInlineSummaryLine(resp.Stdout)
	if inline == "" {
		t.Fatalf("expected inline summary with duration, got:\n%s", resp.Stdout)
	}
	if !strings.Contains(inline, "(3 Run, 2 Pass, 1 Fail, 0 Cached) in ") {
		t.Fatalf("expected (3 Run, 2 Pass, 1 Fail, 0 Cached) in <duration>, got %q", inline)
	}
}
```