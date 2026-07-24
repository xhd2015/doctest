# Scenario

**Feature**: final PASS/FAIL summary surfaces runtime t.Skip counts

```
# fixture leaves include pass / fail / t.Skip
doctest test --no-color <fixture>
  -> go test -json suite-leaf Actions (pass|fail|skip)
  -> when N skip > 0: PASS|FAIL (succeeded/actual_run, N t.Skip) in …
  -> when N = 0: PASS (N/N) in … (unchanged)

# exit: fail=0 → exit 0 even if t.Skip > 0
```

## Preconditions

- The doctest binary is resolved via `testbin.Ensure` from the module root.
- Temp fixture trees are created per leaf (pass / fail / t.Skip combinations).
- Nested CLI self-tests isolate GOCACHE / leaf-cache so concurrent suite noise
  cannot race shared machine caches.
- Outer harness timeout is generous (outer generate + nested `doctest test`).

## Steps

1. Build or reuse the shared doctest binary.
2. Set a generous outer `req.Timeout`.
3. Provide helpers: pass+skip tree, fail+skip tree, all-pass tree, isolated env,
   summary parsers that accept optional `, N t.Skip`.

## Context

- Runtime skip is produced by a fixture leaf calling `t.Skip` (Setup path), so
  go test emits Action `skip` for that suite leaf.
- Product fraction is **succeeded/actual_run** (`actual_run = pass + fail`);
  never planned leaf count when skips exist.
- Primary contract is the final PASS/FAIL line; quiet dots / compact
  `(N Run, …)` line is not asserted here.
- Leaves are labeled `heavy`.

```go
import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")

// finalSummaryPlainRe matches both legacy PASS (p/t) in … and the new form
// PASS (p/t, N t.Skip) in … (and FAIL variants).
var finalSummaryPlainRe = regexp.MustCompile(`^(PASS|FAIL) \(\d+/\d+(?:, \d+ t\.Skip)?\) in .+$`)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Timeout = 180 * time.Second
	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", "..", ".."))
	return nil
}

// isolateRunEnv returns env vars that sandbox go/leaf caches for nested CLI runs.
func isolateRunEnv(t *testing.T) []string {
	t.Helper()
	goCache := t.TempDir()
	leafCache := t.TempDir()
	cacheHome := t.TempDir()
	return []string{
		"GOCACHE=" + goCache,
		"DOCTEST_LEAF_CACHE=" + leafCache,
		"DOCTEST_CACHE_HOME=" + cacheHome,
	}
}

// skipAssertGo is Assert body that always t.Skip (runtime suite-leaf skip).
const skipAssertGo = `import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Skip("intentional runtime t.Skip for summary-tskip fixture")
}
`

// createOnePassOneTSkipTree builds a 2-leaf tree: always-pass + t.Skip.
// Returns the fixture tree root.
func createOnePassOneTSkipTree(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, tmp, []testtree.LeafSpec{
		{
			Name:     "ok",
			Steps:    "always pass",
			Expected: "passes",
		},
		{
			Name:     "skip_me",
			Steps:    "call t.Skip so go test emits Action skip",
			Expected: "runtime skip (not a failure)",
			AssertGo: skipAssertGo,
		},
	})
	return tmp
}

// createOneFailOneTSkipTree builds a 2-leaf tree: forced fail + t.Skip.
func createOneFailOneTSkipTree(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, tmp, []testtree.LeafSpec{
		{
			Name:     "z_fail",
			Steps:    "always fail",
			Expected: "fails",
			AssertGo: `import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Fatal("forced failure")
}
`,
		},
		{
			Name:     "skip_me",
			Steps:    "call t.Skip so go test emits Action skip",
			Expected: "runtime skip (not a failure)",
			AssertGo: skipAssertGo,
		},
	})
	return tmp
}

// createAllPassNoTSkipTree builds a 2-leaf always-pass tree (no t.Skip).
func createAllPassNoTSkipTree(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WritePassFailTree(t, tmp, 2, 0)
	return tmp
}

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

// findResultSummary returns the last PASS/FAIL suite summary line on stdout.
func findResultSummary(stdout string) string {
	var last string
	for _, line := range strings.Split(stdout, "\n") {
		plain := strings.TrimSpace(stripANSI(line))
		if strings.HasPrefix(plain, "PASS (") || strings.HasPrefix(plain, "FAIL (") {
			last = line
		}
	}
	return last
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

func plainSummary(stdout string) string {
	return strings.TrimSpace(stripANSI(findResultSummary(stdout)))
}

// respStdout / respStderr tolerate nil Response in failure messages.
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
