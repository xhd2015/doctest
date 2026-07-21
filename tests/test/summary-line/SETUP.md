# Scenario

**Feature**: end-of-run summary line after `doctest test` completes

```
# discover leaves, run packages, aggregate stats
doctest test <dirs> -> discover leaves -> go test per leaf -> accumulate Passed/Total

# final summary (stdout when Total>0)
runner -> PASS(passed/total) | FAIL(passed/total)

# no runnable leaves
runner -> stderr "no tests" (exit 0)
```

## Preconditions

- The doctest binary is built fresh from module source in this root Setup.
- Temp test trees are created programmatically for each leaf.

## Steps

1. Build a temp doctest tree (or empty dir) matching the leaf scenario.
2. Configure `req.Args` with `doctest test` flags and target paths.
3. Capture stdout/stderr from the subprocess run.

## Context

- Non-verbose runs include per-suite `(N Run, N Pass, N Fail, N Cached)` after dots.
- The aggregated `PASS(x/y)` or `FAIL(x/y)` line is the last non-empty stdout line when cases exist.
- Color tests use `--color` or `--no-color` CLI flags.
- Multi-tree prepare failures use `createPrepareFailMultiTree` + `./...` from module root.

```go
import (
"github.com/xhd2015/doctest/session"
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
var finalSummaryPlainRe = regexp.MustCompile(`^(PASS|FAIL) \(\d+/\d+\) in .+$`)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Timeout = 120 * time.Second

	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", "..", ".."))
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

// createPrepareFailMultiTree builds one module with a good 1-pass tree and a
// sibling tree whose DOCTEST.md has invalid Go so prepare fails. Caller runs
// from the returned module root with `doctest test ./...`.
func createPrepareFailMultiTree(t *testing.T) string {
	t.Helper()
	mod := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(mod, "go.mod"),
		[]byte("module example.com/prepfail\n\ngo 1.22\n"),
		0644,
	); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	good := filepath.Join(mod, "good")
	testtree.WritePassFailTree(t, good, 1, 0)

	bad := filepath.Join(mod, "bad")
	// Invalid Go in DOCTEST.md so generate/prepare fails (syntax error).
	testtree.WriteFile(t, bad, "DOCTEST.md", testtree.MinimalDOCTEST(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil
// missing closing brace — syntax error
`))
	fence := bt(3)
	setupBody := `import "testing"

func Setup(t *testing.T, req *Request) error { _ = req; return nil }`
	assertBody := `func Assert(t *testing.T, req *Request, resp *Response, err error) {}`
	testtree.WriteFile(t, bad, "leaf/SETUP.md",
		"## Steps\n1. bad leaf\n\n"+fence+"go\n"+setupBody+"\n"+fence+"\n")
	testtree.WriteFile(t, bad, "leaf/ASSERT.md",
		"## Expected\n- n/a\n\n"+fence+"go\n"+assertBody+"\n"+fence+"\n")
	return mod
}

func createEmptyDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
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
func findInlineSummaryLine(stdout string) string {
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, " Run, ") && strings.Contains(line, " Cached") {
			return line
		}
	}
	return ""
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
