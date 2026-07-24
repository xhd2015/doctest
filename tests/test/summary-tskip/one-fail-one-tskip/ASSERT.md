---
label: heavy
explanation: nested doctest test on fail+t.Skip fixture; asserts FAIL (0/1, 1 t.Skip) and non-zero exit
---

## Expected

- Process exits **non-zero** (at least one real fail).
- Final suite summary contains **`FAIL (`** and **`t.Skip`**.
- Fraction is **succeeded/actual_run** = **0/1** (one fail, skip excluded).
- Target form: **`FAIL (0/1, 1 t.Skip) in <duration>`**.
- Must **not** use planned 2 in the denominator (`0/2` or `1/2`).
- Summary is the last non-empty stdout line.

## Side Effects

- Temp fixture tree under the leaf TempDir.
- Isolated GOCACHE / leaf-cache for the harness invocation.

## Errors

- Treating runtime skip as a pass or fail in the fraction.
- Using planned leaf count as denominator.

## Exit Code

- Non-zero.

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
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit (real fail present), got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	plain := plainSummary(resp.Stdout)
	if plain == "" {
		t.Fatalf("expected PASS/FAIL summary line, got:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.HasPrefix(plain, "FAIL (") {
		t.Fatalf("expected FAIL (… summary, got %q\nstdout:\n%s", plain, resp.Stdout)
	}
	if !strings.Contains(plain, "t.Skip") {
		t.Fatalf("expected t.Skip on FAIL summary when runtime skip occurred, got %q\nstdout:\n%s",
			plain, resp.Stdout)
	}
	// Reject planned-denominator forms.
	if strings.Contains(plain, "FAIL (0/2") || strings.Contains(plain, "FAIL (1/2") {
		t.Fatalf("must not use planned leaf count in denominator; got %q\nstdout:\n%s",
			plain, resp.Stdout)
	}
	if !strings.HasPrefix(plain, "FAIL (0/1, 1 t.Skip) in ") {
		t.Fatalf("expected FAIL (0/1, 1 t.Skip) in <duration>, got %q\nstdout:\n%s\nstderr:\n%s",
			plain, resp.Stdout, resp.Stderr)
	}
	if _, err := parseFinalSummaryDuration(resp.Stdout); err != nil {
		t.Fatalf("final duration must parse: %v\nsummary: %q", err, plain)
	}
	summaryIsLastLine(t, resp.Stdout)
}
```
