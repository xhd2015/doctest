# Scenario

**Feature**: when nested go test times out, doctest surfaces Error/hint, FAIL cancelled accounting, and optional color accents

```
# user passes short timeout; multi-leaf suite sleeps past it
doctest test --timeout=2s <sleep-tree> -> go test -timeout=2s -> panic: test timed out after 2s

# fail path: Error + hint, then FAIL with planned denom + cancelled
doctest -> Error: go test timed out after 2s
doctest -> hint: increase with -timeout=DURATION (e.g. -timeout=30m; -timeout=0 disables)
doctest -> FAIL (0/3, N cancelled) in …   # N = planned − pass − fail − skip > 0

# color on: Error red, hint gray, "N cancelled" orange
# progress stays finished-only (no Cancelled segment)
# fast pass must not emit timeout Error or cancelled
```

## Preconditions

- The doctest binary is resolved via `testbin.Ensure` from the module root.
- Temp fixture trees are created per leaf (multi-sleep or fast-pass).
- Outer harness timeout is generous (compile + nested 2s timeout).

## Steps

1. Build or reuse the shared doctest binary.
2. Set a generous outer `req.Timeout` (nested runs compile + may wait on timeout).
3. Provide helpers to write multi-leaf sleep and fast-pass fixture trees, plus
   output parsers for timeout Error/hint, FAIL cancelled summary, print order,
   and ANSI color checks (red / gray / orange `38;5;208`).

## Context

- Nested CLI self-tests use `label: e2e` when full-integration.
- Default timeout policy is unchanged: only messaging when a timeout actually fires.
- Locked Error/hint wording; FAIL uses planned denom when cancelled > 0.
- Progress compact line must not gain a Cancelled segment.
- Preferred color: orange only on the cancelled phrase (warning), not on Error.

```go
import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

const (
	ansiReset  = "\x1b[0m"
	ansiRed    = "\x1b[31m"
	ansiGray   = "\x1b[90m"
	ansiOrange = "\x1b[38;5;208m" // warning accent for "N cancelled"
)

var (
	ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")
	// FAIL (p/P, N cancelled) in <duration>
	failCancelledPlainRe = regexp.MustCompile(
		`^FAIL \((\d+)/(\d+), (\d+) cancelled\) in .+$`,
	)
	timeoutErrorPlainRe = regexp.MustCompile(
		`^Error: go test timed out after \S+$`,
	)
)

const lockedHint = "hint: increase with -timeout=DURATION (e.g. -timeout=30m; -timeout=0 disables)"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Timeout = 180 * time.Second
	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", "..", ".."))
	return nil
}

// createSleepTree builds a multi-leaf fixture whose shared Run sleeps for
// sleepSec seconds so nested go test -timeout can fire. leafCount is the
// discovery planned total used for cancelled accounting.
func createSleepTree(t *testing.T, leafCount, sleepSec int) string {
	t.Helper()
	if leafCount < 1 {
		t.Fatalf("leafCount must be >= 1, got %d", leafCount)
	}
	tmp := t.TempDir()
	leaves := make([]testtree.LeafSpec, 0, leafCount)
	for i := 0; i < leafCount; i++ {
		leaves = append(leaves, testtree.LeafSpec{
			Name:     fmt.Sprintf("sleep_%d", i),
			Steps:    "Run sleeps past go test -timeout",
			Expected: "would pass if not timed out",
		})
	}
	testtree.WriteMinimalRunnableTree(t, tmp, leaves)
	runGo := fmt.Sprintf(`import (
	"testing"
	"time"
)

type Request struct{}
type Response struct{}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	time.Sleep(%d * time.Second)
	return &Response{}, nil
}`, sleepSec)
	// Overwrite DOCTEST.md Run so the nested suite actually sleeps.
	testtree.WriteFile(t, tmp, "DOCTEST.md", testtree.MinimalDOCTEST(runGo))
	return tmp
}

// createFastPassTree builds a 1-pass fixture that finishes well under any normal timeout.
func createFastPassTree(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WritePassFailTree(t, tmp, 1, 0)
	return tmp
}

// combinedOutput joins stdout and stderr for fail-path visibility checks.
func combinedOutput(resp *Response) string {
	if resp == nil {
		return ""
	}
	return resp.Stdout + "\n" + resp.Stderr
}

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

func containsANSI(s string) bool {
	return ansiEscape.MatchString(s)
}

// hasTimeoutError reports the locked Error line (plain text, after strip).
func hasTimeoutError(combined string) bool {
	for _, line := range strings.Split(stripANSI(combined), "\n") {
		if timeoutErrorPlainRe.MatchString(strings.TrimSpace(line)) {
			return true
		}
	}
	// also accept the exact preferred form substring
	return strings.Contains(stripANSI(combined), "Error: go test timed out after")
}

// hasLockedHint reports the locked hint line.
func hasLockedHint(combined string) bool {
	return strings.Contains(stripANSI(combined), lockedHint)
}

// findResultSummary returns the first PASS/FAIL summary line (may include ANSI).
func findResultSummary(s string) string {
	for _, line := range strings.Split(s, "\n") {
		plain := strings.TrimSpace(stripANSI(line))
		if strings.HasPrefix(plain, "PASS (") || strings.HasPrefix(plain, "FAIL (") {
			return line
		}
	}
	return ""
}

// findInlineProgressSummary returns the quiet compact progress summary line.
func findInlineProgressSummary(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(line, " Run, ") && strings.Contains(line, " Cached") {
			return line
		}
	}
	return ""
}

// parseFailCancelled extracts passed, planned, cancelled from a FAIL summary.
// Returns ok=false when the cancelled form is absent.
func parseFailCancelled(summaryLine string) (passed, planned, cancelled int, ok bool) {
	plain := strings.TrimSpace(stripANSI(summaryLine))
	m := failCancelledPlainRe.FindStringSubmatch(plain)
	if m == nil {
		return 0, 0, 0, false
	}
	passed, _ = strconv.Atoi(m[1])
	planned, _ = strconv.Atoi(m[2])
	cancelled, _ = strconv.Atoi(m[3])
	return passed, planned, cancelled, true
}

// timeoutErrorBeforeFail reports whether the locked Error line appears before
// the FAIL (… summary on the primary user-facing stream (stdout). Classic TDD:
// Error/hint must be printed before the final FAIL line.
func timeoutErrorBeforeFail(stdout string) bool {
	plain := stripANSI(stdout)
	errIdx := strings.Index(plain, "Error: go test timed out after")
	failIdx := -1
	for _, line := range strings.Split(plain, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "FAIL (") {
			failIdx = strings.Index(plain, line)
			if failIdx < 0 {
				failIdx = strings.Index(plain, trim)
			}
			break
		}
	}
	if errIdx < 0 || failIdx < 0 {
		return false
	}
	return errIdx < failIdx
}

// lineContaining returns the first raw line whose plain text contains sub.
func lineContaining(s, sub string) string {
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(stripANSI(line), sub) {
			return line
		}
	}
	return ""
}

// segmentColored reports whether needle appears wrapped in openSGR…reset in s.
// Accepts openSGR immediately before needle (optionally after other SGRs on the line).
func segmentColored(s, openSGR, needle string) bool {
	if !strings.Contains(s, openSGR) || !strings.Contains(stripANSI(s), needle) {
		return false
	}
	// Prefer exact open + needle + reset.
	if strings.Contains(s, openSGR+needle+ansiReset) {
		return true
	}
	// Nested / re-open forms: openSGR then needle before next hard break.
	idx := strings.Index(s, needle)
	if idx < 0 {
		return false
	}
	before := s[:idx]
	return strings.Contains(before, openSGR)
}
```
