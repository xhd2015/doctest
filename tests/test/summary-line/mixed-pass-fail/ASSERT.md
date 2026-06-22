## Expected

- Exit code is non-zero (at least one failure).
- stdout contains `FAIL\t` lines before the summary.
- Final stdout line is `FAIL (2/3)`.
- stdout includes inline `(3 Run, 2 Pass, 1 Fail, 0 Cached)` after dots.

## Exit Code

- Non-zero.

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
	if !strings.Contains(resp.Stdout, "...  (3 Run, 2 Pass, 1 Fail, 0 Cached)") {
		t.Fatalf("expected inline progress summary, got:\n%s", resp.Stdout)
	}

	hasFailTab := false
	summaryIdx := -1
	for i, line := range strings.Split(resp.Stdout, "\n") {
		if strings.HasPrefix(line, "FAIL\t") {
			hasFailTab = true
		}
		plain := strings.TrimSpace(stripANSI(line))
		if strings.HasPrefix(plain, "FAIL (") || strings.HasPrefix(plain, "PASS (") {
			summaryIdx = i
		}
	}
	if !hasFailTab {
		t.Fatalf("expected FAIL\\t lines before summary, got:\n%s", resp.Stdout)
	}

	summary := findResultSummary(resp.Stdout)
	if stripANSI(strings.TrimSpace(summary)) != "FAIL (2/3)" {
		t.Fatalf("expected FAIL (2/3) summary, got %q\nstdout:\n%s", summary, resp.Stdout)
	}
	summaryIsLastLine(t, resp.Stdout)

	if summaryIdx >= 0 {
		before := strings.Join(strings.Split(resp.Stdout, "\n")[:summaryIdx], "\n")
		if !strings.Contains(before, "FAIL\t") {
			t.Fatalf("FAIL\\t lines must appear before summary line\nstdout:\n%s", resp.Stdout)
		}
	}
}
```