---
label: heavy
explanation: regression — zero runtime skips must keep PASS (2/2) with no t.Skip text
---

## Expected

- Process exits **0**.
- Final suite summary is **`PASS (2/2) in <duration>`**.
- Summary must **not** contain `t.Skip` (N=0 keeps legacy form).
- Summary is the last non-empty stdout line.

## Side Effects

- Temp 2-pass fixture under the leaf TempDir.
- Isolated GOCACHE / leaf-cache for the harness invocation.

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
		t.Fatalf("expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	plain := plainSummary(resp.Stdout)
	if plain == "" {
		t.Fatalf("expected PASS/FAIL summary line, got:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if strings.Contains(plain, "t.Skip") {
		t.Fatalf("N=0 must not include t.Skip text on summary, got %q\nstdout:\n%s",
			plain, resp.Stdout)
	}
	if !strings.HasPrefix(plain, "PASS (2/2) in ") {
		t.Fatalf("expected PASS (2/2) in <duration> with no t.Skip, got %q\nstdout:\n%s\nstderr:\n%s",
			plain, resp.Stdout, resp.Stderr)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v\nsummary: %q", err, plain)
	}
	summaryIsLastLine(t, resp.Stdout)
}
```
