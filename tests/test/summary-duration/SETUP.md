# Scenario

**Feature**: elapsed time in per-suite and final test summaries

```
# discover leaves, run packages, measure wall time
doctest test <dirs> -> discover leaves -> go test per suite -> inline (N Run, ...) in DURATION

# final aggregate
runner -> PASS(passed/total) in DURATION | FAIL(passed/total) in DURATION
```

## Preconditions

- The doctest binary is built fresh from module source in this root Setup.
- Temp test trees are created programmatically for each leaf.

## Steps

1. Build a temp doctest tree (or empty dir) matching the leaf scenario.
2. Configure `req.Args` with `doctest test` flags and target paths.
3. Capture stdout/stderr from the subprocess run.

## Context

- Non-verbose runs include per-suite `(N Run, N Pass, N Fail, N Cached) in DURATION` after dots.
- The aggregated `PASS (x/y) in DURATION` or `FAIL (x/y) in DURATION` line is the last non-empty stdout line when cases exist.
- Color tests use `--color` or `--no-color` CLI flags.
- Duration uses display formatting: sub-second values as integers (e.g. `949ms`); ≥1s with at most 2 decimal digits (e.g. `1.37s`). Assertions parse rather than match exact values.

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")
var inlineSummaryPlainRe = regexp.MustCompile(`\(\d+ Run, \d+ Pass, \d+ Fail, \d+ Cached\) in (.+)`)
var finalSummaryPlainRe = regexp.MustCompile(`^(PASS|FAIL) \(\d+/\d+\) in .+$`)
func Setup(t *testing.T, req *Request) error {
	req.Timeout = 120 * time.Second

	req.Bin = testbin.Ensure(t, filepath.Join(DOCTEST_ROOT, "..", "..", ".."))
	return nil
}
func bt(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = '`'
	}
	return string(b)
}
func doctestGoBlock(code string) string {
	fence := bt(3)
	return "## Test\n\n" + fence + "go\n" + code + "\n" + fence + "\n"
}
func createPassFailTree(t *testing.T, passCount int, failCount int) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WritePassFailTree(t, tmp, passCount, failCount)
	return tmp
}
func createSlowLeafTree(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WriteFile(t, tmp, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))

	leafDir := filepath.Join(tmp, "slow_1")
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		t.Fatal(err)
	}
	slowSetup := `import (
	"testing"
	"time"
)
func Setup(t *testing.T, req *Request) error {
	time.Sleep(time.Second)
	return nil
}`
	if err := os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(doctestGoBlock(slowSetup)), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(doctestGoBlock(`import "testing"
func Assert(t *testing.T, req *Request, resp *Response, err error) {}`)), 0644); err != nil {
		t.Fatal(err)
	}

	return tmp
}
func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}
func containsANSI(s string) bool {
	return ansiEscape.MatchString(s)
}
func lastNonEmptyLine(s string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			return lines[i]
		}
	}
	return ""
}
func findResultSummary(stdout string) string {
	for _, line := range strings.Split(stdout, "\n") {
		plain := strings.TrimSpace(stripANSI(line))
		if strings.HasPrefix(plain, "PASS (") || strings.HasPrefix(plain, "FAIL (") {
			return line
		}
	}
	return ""
}
func countResultSummaries(stdout string) int {
	n := 0
	for _, line := range strings.Split(stdout, "\n") {
		plain := strings.TrimSpace(stripANSI(line))
		if strings.HasPrefix(plain, "PASS (") || strings.HasPrefix(plain, "FAIL (") {
			n++
		}
	}
	return n
}
func findInlineSummaryLine(stdout string) string {
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, " Run, ") && strings.Contains(line, " Cached") {
			return line
		}
	}
	return ""
}
func countInlineSummaries(stdout string) int {
	n := 0
	for _, line := range strings.Split(stdout, "\n") {
		plain := stripANSI(strings.TrimSpace(line))
		if inlineSummaryPlainRe.MatchString(plain) {
			n++
		}
	}
	return n
}
func parseInlineSummaryDuration(line string) (time.Duration, error) {
	plain := stripANSI(strings.TrimSpace(line))
	matches := inlineSummaryPlainRe.FindStringSubmatch(plain)
	if len(matches) < 2 {
		return 0, fmt.Errorf("inline summary not found in %q", line)
	}
	return time.ParseDuration(strings.TrimSpace(matches[1]))
}
func parseFinalSummaryDuration(stdout string) (time.Duration, error) {
	line := stripANSI(strings.TrimSpace(findResultSummary(stdout)))
	if line == "" {
		return 0, fmt.Errorf("final summary line not found")
	}
	if !finalSummaryPlainRe.MatchString(line) {
		return 0, fmt.Errorf("final summary missing duration suffix: %q", line)
	}
	const marker = " in "
	idx := strings.LastIndex(line, marker)
	if idx < 0 {
		return 0, fmt.Errorf("final summary missing %q: %q", marker, line)
	}
	return time.ParseDuration(strings.TrimSpace(line[idx+len(marker):]))
}
func inlineDurationIsGray(stdout string) bool {
	line := findInlineSummaryLine(stdout)
	if line == "" {
		return false
	}
	const marker = ") in "
	idx := strings.LastIndex(line, marker)
	if idx < 0 {
		return false
	}
	durSeg := line[idx+len(marker):]
	return containsANSI(durSeg)
}
func finalSummaryPassTokenIsColored(summary string) bool {
	idx := strings.Index(summary, " in ")
	if idx < 0 {
		return false
	}
	return containsANSI(summary[:idx])
}
func finalSummaryDurationIsPlain(summary string) bool {
	idx := strings.LastIndex(summary, " in ")
	if idx < 0 {
		return false
	}
	return !containsANSI(summary[idx:])
}
func summaryIsLastLine(t *testing.T, stdout string) {
	t.Helper()
	summary := findResultSummary(stdout)
	if summary == "" {
		t.Fatalf("expected PASS/FAIL summary line in stdout:\n%s", stdout)
	}
	last := lastNonEmptyLine(stdout)
	if strings.TrimSpace(last) != strings.TrimSpace(summary) {
		t.Fatalf("summary must be last non-empty stdout line\nsummary: %q\nlast: %q\nstdout:\n%s",
			summary, last, stdout)
	}
}
```
