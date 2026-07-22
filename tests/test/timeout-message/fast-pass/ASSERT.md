---
label: heavy
explanation: nested doctest test compiles and runs a fast 1-pass fixture
---

## Expected

- Command succeeds (exit 0).
- Combined stdout+stderr does **not** contain `Error: go test timed out`
  (no false positive timeout messaging).
- Combined output does **not** contain `cancelled` (no cancelled phrase when
  there was no timeout).
- Final summary is ordinary PASS without a cancelled segment.

## Side Effects

- None beyond subprocess output and temp fixture files.

## Exit Code

- 0.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("harness run failed unexpectedly: %v\nstdout:\n%s\nstderr:\n%s",
			err, resp.Stdout, resp.Stderr)
	}
	if resp == nil {
		t.Fatal("expected non-nil Response")
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0 for fast-pass tree, got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	combined := combinedOutput(resp)
	if strings.Contains(combined, "Error: go test timed out") {
		t.Fatalf("false positive: found \"Error: go test timed out\" on a passing suite\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if strings.Contains(strings.ToLower(stripANSI(combined)), "cancelled") {
		t.Fatalf("false positive: found \"cancelled\" on a passing suite without timeout\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	summary := findResultSummary(resp.Stdout)
	if summary != "" {
		plain := strings.TrimSpace(stripANSI(summary))
		if strings.Contains(plain, "cancelled") {
			t.Fatalf("PASS/FAIL summary must not include cancelled without timeout: %q", plain)
		}
	}
}
```
