---
label: heavy
explanation: outer 2-leaf tree plus nested doctest test on intentional fail; -v stream must still report PASS (2/2)
---

## Expected

- Both outer leaves pass (including `nested_fail_ok`, which only *hosts* a nested fail).
- Nested stdout is leaked into the harness stream: combined output contains `FAIL (`.
- Process exits 0.
- Final suite summary is **`PASS (2/2) in <duration>`** — not `FAIL (1/2)`,
  `FAIL (0/2)`, or any `FAIL (p/2)` with p < 2.
- Summary is the last non-empty stdout line.

## Side Effects

- Nested intentional-fail child and outer temp trees under the leaf TempDir.
- Isolated GOCACHE / leaf-cache for harness and nested invocations.

## Errors

- Must **not** treat nested `FAIL (` / `FAIL\t` text as outer package failures
  when counting suite Passed under `-v`.

## Exit Code

- 0.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("harness run failed unexpectedly: %v\nstdout:\n%s\nstderr:\n%s",
			err, respStdout(resp), respStderr(resp))
	}
	if resp == nil {
		t.Fatal("expected non-nil Response")
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0 (both outer leaves pass), got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	combined := combinedOutput(resp)
	if !streamHasNestedFailSummary(combined) {
		t.Fatalf("precondition: expected nested FAIL ( in stream (proof nested fail leaked), got:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	summary := findResultSummary(resp.Stdout)
	plain := strings.TrimSpace(stripANSI(summary))
	if !strings.HasPrefix(plain, "PASS (2/2) in ") {
		t.Fatalf("expected PASS (2/2) in <duration> under -v despite nested FAIL ( text, got %q\nstdout:\n%s\nstderr:\n%s",
			plain, resp.Stdout, resp.Stderr)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v\nsummary: %q", err, summary)
	}
	summaryIsLastLine(t, resp.Stdout)
}

func respStdout(resp *Response) string {
	if resp == nil {
		return ""
	}
	return resp.Stdout
}

func respStderr(resp *Response) string {
	if resp == nil {
		return ""
	}
	return resp.Stderr
}
```
