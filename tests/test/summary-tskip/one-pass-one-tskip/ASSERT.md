---
label: heavy
explanation: nested doctest test on 2-leaf pass+t.Skip fixture; asserts new PASS (1/1, 1 t.Skip) form
---

## Expected

- Process exits **0** (runtime skip alone must not fail the suite).
- Final suite summary is **`PASS (1/1, 1 t.Skip) in <duration>`**.
- Fraction is **succeeded/actual_run** = 1/1 (skips excluded from denominator).
- Must **not** be bare `PASS (1/1)` without `t.Skip`.
- Must **not** be `PASS (1/2, …)` (planned total wrongly used as denominator).
- Summary is the last non-empty stdout line.

## Side Effects

- Temp fixture tree under the leaf TempDir.
- Isolated GOCACHE / leaf-cache for the harness invocation.

## Errors

- Wrong denominator using planned leaf count (`1/2`).
- Missing `, N t.Skip` suffix when a runtime skip occurred.

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
		t.Fatalf("expected exit 0 (pass + t.Skip must not fail suite), got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	plain := plainSummary(resp.Stdout)
	if plain == "" {
		t.Fatalf("expected PASS/FAIL summary line, got:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	// Reject wrong planned-denominator form first (clearer diagnostic).
	if strings.Contains(plain, "PASS (1/2") || strings.Contains(plain, "FAIL (1/2") {
		t.Fatalf("must not use planned leaf count in denominator; got %q\nstdout:\n%s",
			plain, resp.Stdout)
	}
	if !strings.HasPrefix(plain, "PASS (1/1, 1 t.Skip) in ") {
		t.Fatalf("expected PASS (1/1, 1 t.Skip) in <duration>, got %q\nstdout:\n%s\nstderr:\n%s",
			plain, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(plain, "t.Skip") {
		t.Fatalf("expected t.Skip suffix on summary when runtime skip occurred, got %q", plain)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v\nsummary: %q", err, plain)
	}
	summaryIsLastLine(t, resp.Stdout)
}
```
