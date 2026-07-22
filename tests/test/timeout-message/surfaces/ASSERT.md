---
label: heavy
explanation: nested doctest test compiles a 3-leaf sleep fixture and waits on go test -timeout=2s
---

## Expected

- Nested `go test` times out (sleep ≥ timeout).
- Process exits non-zero.
- Combined stdout+stderr contains locked timeout wording:

  ```
  Error: go test timed out after 2s
  hint: increase with -timeout=DURATION (e.g. -timeout=30m; -timeout=0 disables)
  ```

- Final suite summary (plain) is:

  `FAIL (0/3, N cancelled) in <duration>`

  with **planned = 3** (discovery leaf count), **N > 0**, and N consistent with
  `cancelled = planned − pass − fail − skip` (v1: no `t.Skip` phrase on this line).
- Quiet progress compact line has **no** `cancelled` / `Cancelled` segment
  (finished-only Run/Pass/Fail/Cached).
- Print order on **stdout**: timeout Error (and hint) appear **before** the
  final `FAIL (` summary line.

## Errors

- Timeout is reported as a first-class fail-path message, not only buried
  under filtered JSON / progress dots.
- Must not keep legacy timeout FAIL as bare `FAIL (0/1)` (package actual_run)
  when planned leaves were cancelled.

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
		t.Fatalf("expected non-zero exit when nested go test times out, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	combined := combinedOutput(resp)
	if !hasTimeoutError(combined) {
		t.Fatalf("expected locked Error line \"Error: go test timed out after …\" in stdout+stderr, got exit=%d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !hasLockedHint(combined) {
		t.Fatalf("expected locked hint line %q in stdout+stderr\nstdout:\n%s\nstderr:\n%s",
			lockedHint, resp.Stdout, resp.Stderr)
	}

	// Progress: finished-only; no cancelled segment.
	progress := findInlineProgressSummary(resp.Stdout)
	if progress != "" {
		plainProg := stripANSI(progress)
		if strings.Contains(strings.ToLower(plainProg), "cancelled") {
			t.Fatalf("progress line must not include cancelled; got %q", plainProg)
		}
	}

	// Prefer summary from stdout (user-facing end line).
	summary := findResultSummary(resp.Stdout)
	if summary == "" {
		summary = findResultSummary(combined)
	}
	if summary == "" {
		t.Fatalf("expected FAIL summary line, got:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	plainSummary := strings.TrimSpace(stripANSI(summary))
	if strings.Contains(plainSummary, "t.Skip") {
		t.Fatalf("v1 timeout FAIL must omit t.Skip phrase; got %q", plainSummary)
	}
	passed, planned, cancelled, ok := parseFailCancelled(summary)
	if !ok {
		t.Fatalf("expected FAIL (0/3, N cancelled) in <duration>, got %q\nstdout:\n%s\nstderr:\n%s",
			plainSummary, resp.Stdout, resp.Stderr)
	}
	if planned != 3 {
		t.Fatalf("FAIL denom planned must be discovery leaf count 3, got planned=%d summary=%q",
			planned, plainSummary)
	}
	if cancelled <= 0 {
		t.Fatalf("expected cancelled > 0 on multi-leaf timeout, got %d summary=%q",
			cancelled, plainSummary)
	}
	if passed < 0 || passed > planned {
		t.Fatalf("unexpected passed=%d planned=%d summary=%q", passed, planned, plainSummary)
	}
	// Accounting: cancelled ≤ planned − passed (fail/skip non-negative residual).
	if cancelled > planned-passed {
		t.Fatalf("cancelled %d exceeds planned-passed %d; summary=%q",
			cancelled, planned-passed, plainSummary)
	}

	// Error/hint before FAIL on stdout (print order).
	if !timeoutErrorBeforeFail(resp.Stdout) {
		// If Error only on stderr today, still RED: product must put Error before FAIL on stdout.
		t.Fatalf("expected timeout Error (and hint) on stdout before FAIL summary\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if idxErr := strings.Index(stripANSI(resp.Stdout), "Error: go test timed out after"); idxErr >= 0 {
		if idxHint := strings.Index(stripANSI(resp.Stdout), lockedHint); idxHint >= 0 && idxHint < idxErr {
			t.Fatalf("hint must not appear before Error line\nstdout:\n%s", resp.Stdout)
		}
	}

	// --no-color: summary / Error lines must be plain (no ANSI on those lines).
	if errLine := lineContaining(resp.Stdout, "Error: go test timed out after"); errLine != "" && containsANSI(errLine) {
		t.Fatalf("--no-color Error line must not contain ANSI, got %q", errLine)
	}
	if containsANSI(summary) {
		t.Fatalf("--no-color FAIL summary must not contain ANSI, got %q", summary)
	}
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
