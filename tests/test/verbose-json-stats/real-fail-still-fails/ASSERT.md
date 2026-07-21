---
label: heavy
explanation: regression guard — real outer Assert fail still yields FAIL (0/1) and non-zero exit
---

## Expected

- The single outer leaf fails.
- Process exits non-zero.
- Final suite summary is **`FAIL (0/1) in <duration>`** (p < t).
- Summary is the last non-empty stdout line.

## Side Effects

- Temp 1-fail fixture under the leaf TempDir.
- Isolated GOCACHE / leaf-cache for the harness invocation.

## Errors

- Real failures remain first-class: always-json must not mask genuine fails.

## Exit Code

- Non-zero.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("harness run failed unexpectedly: %v\nstdout:\n%s\nstderr:\n%s",
			err, respStdout(resp), respStderr(resp))
	}
	if resp == nil {
		t.Fatal("expected non-nil Response")
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit for real outer fail, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	summary := findResultSummary(resp.Stdout)
	plain := strings.TrimSpace(stripANSI(summary))
	if !strings.HasPrefix(plain, "FAIL (0/1) in ") {
		t.Fatalf("expected FAIL (0/1) in <duration> for real outer fail, got %q\nstdout:\n%s\nstderr:\n%s",
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
