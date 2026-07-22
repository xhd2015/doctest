---
label: heavy
explanation: nested doctest test with --color on 3-leaf sleep fixture; timeout colors + orange cancelled
---

## Expected

- Nested `go test` times out; exit non-zero.
- Same plain wording as `surfaces` after stripANSI:
  - `Error: go test timed out after 2s`
  - locked hint line
  - `FAIL (0/3, N cancelled) in <duration>` with N > 0
- When color is on:
  - **Error** line uses **red** (`\x1b[31m`)
  - **hint** line uses **gray** (`\x1b[90m`)
  - **`N cancelled`** segment uses **orange** (`\x1b[38;5;208m`) — warning accent
  - FAIL token remains red (nested orange on cancelled is OK)
- Progress line still has no `cancelled` segment.
- Error appears on stdout before FAIL summary (same order contract as surfaces).

## Errors

- Missing orange on cancelled while Error is red (wrong accent assignment).
- Coloring Error orange, or leaving cancelled plain under `--color`.

## Exit Code

- Non-zero.

```go
import (
	"fmt"
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
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit when nested go test times out, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}

	// Prefer stdout for colored user-facing lines (print order + color).
	out := resp.Stdout
	combined := combinedOutput(resp)

	if !hasTimeoutError(combined) {
		t.Fatalf("expected Error: go test timed out after …\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !hasLockedHint(combined) {
		t.Fatalf("expected locked hint %q\nstdout:\n%s\nstderr:\n%s",
			lockedHint, resp.Stdout, resp.Stderr)
	}

	errLine := lineContaining(out, "Error: go test timed out after")
	if errLine == "" {
		// Color + order require Error on stdout under --color.
		t.Fatalf("expected timeout Error on stdout under --color\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !segmentColored(errLine, ansiRed, "Error: go test timed out after") &&
		!strings.Contains(errLine, ansiRed) {
		t.Fatalf("expected red ANSI on Error line, got %q", errLine)
	}

	hintLine := lineContaining(out, lockedHint)
	if hintLine == "" {
		hintLine = lineContaining(combined, lockedHint)
	}
	if hintLine == "" {
		t.Fatalf("expected hint line in output\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	// Hint must be gray when color on (and not red-as-Error).
	if !segmentColored(hintLine, ansiGray, "hint:") && !strings.Contains(hintLine, ansiGray) {
		t.Fatalf("expected gray ANSI on hint line, got %q", hintLine)
	}
	if strings.Contains(hintLine, ansiRed) && !strings.Contains(hintLine, ansiGray) {
		t.Fatalf("hint must not be red-only; got %q", hintLine)
	}

	summary := findResultSummary(out)
	if summary == "" {
		summary = findResultSummary(combined)
	}
	if summary == "" {
		t.Fatalf("expected FAIL summary\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	passed, planned, cancelled, ok := parseFailCancelled(summary)
	if !ok || planned != 3 || cancelled <= 0 {
		t.Fatalf("expected FAIL (0/3, N cancelled) with N>0, got plain %q (passed=%d planned=%d cancelled=%d ok=%v)",
			strings.TrimSpace(stripANSI(summary)), passed, planned, cancelled, ok)
	}
	cancelledPhrase := fmt.Sprintf("%d cancelled", cancelled)
	if !segmentColored(summary, ansiOrange, cancelledPhrase) &&
		!strings.Contains(summary, ansiOrange) {
		t.Fatalf("expected orange (%s) on %q in FAIL summary, got %q",
			strings.TrimPrefix(ansiOrange, "\x1b["), cancelledPhrase, summary)
	}
	// FAIL token red (or at least red present on the summary line).
	if !strings.Contains(summary, ansiRed) {
		t.Fatalf("expected red ANSI on FAIL summary token, got %q", summary)
	}
	// Orange is for cancelled only — Error stays red, not orange.
	if strings.Contains(errLine, ansiOrange) {
		t.Fatalf("Error line must stay red (not orange); got %q", errLine)
	}

	progress := findInlineProgressSummary(out)
	if progress != "" && strings.Contains(strings.ToLower(stripANSI(progress)), "cancelled") {
		t.Fatalf("progress must not include cancelled; got %q", progress)
	}
	if !timeoutErrorBeforeFail(out) {
		t.Fatalf("expected Error before FAIL on stdout under --color\nstdout:\n%s", out)
	}
}
```
