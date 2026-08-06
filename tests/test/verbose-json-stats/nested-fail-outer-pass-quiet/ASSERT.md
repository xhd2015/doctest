---
explanation: same nested-fail outer-pass fixture without -v; json path should already report PASS (2/2)
---

## Expected

- Both outer leaves pass.
- Process exits 0.
- Final suite summary is **`PASS (2/2) in <duration>`**.
- Same Passed/Total as the verbose sibling (quiet is the control for always-json).

## Side Effects

- Nested intentional-fail child and outer temp trees under the leaf TempDir.
- Isolated GOCACHE / leaf-cache for harness and nested invocations.

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
	summary := findResultSummary(resp.Stdout)
	plain := strings.TrimSpace(stripANSI(summary))
	if !strings.HasPrefix(plain, "PASS (2/2) in ") {
		t.Fatalf("expected PASS (2/2) in <duration> on quiet path, got %q\nstdout:\n%s\nstderr:\n%s",
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
