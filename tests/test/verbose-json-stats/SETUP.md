# Scenario

**Feature**: suite Pass/Fail counts always come from go test -json, even under -v

```
# outer suite: 2 leaves both Assert-pass; nested intentional fail leaks FAIL (
doctest test -v --no-color <outer>
  -> go test (stats via -json, not text FAIL scan)
  -> stream may contain nested "FAIL (0/1)"
  -> final summary PASS (2/2)

# quiet same tree
doctest test --no-color <outer> -> PASS (2/2)

# real outer fail still fails
doctest test --no-color <1-fail> -> FAIL (0/1), exit ≠ 0
```

## Preconditions

- The doctest binary is resolved via `testbin.Ensure` from the module root.
- Temp fixture trees are created per leaf (outer + optional nested child).
- Nested CLI self-tests isolate GOCACHE / leaf-cache so concurrent suite noise
  cannot race shared machine caches.
- Outer harness timeout is generous (outer generate + nested `doctest test`).

## Steps

1. Build or reuse the shared doctest binary.
2. Set a generous outer `req.Timeout`.
3. Provide helpers: nested-fail outer-pass fixture, real-fail fixture, isolated
   env, summary parsers.

## Context

- Nested intentional fail is **not** an outer leaf failure: outer Assert expects
  nested non-zero and prints nested stdout so `-v` streams contain `FAIL (`.
- Until always-json lands, text scanning of `FAIL (` / `FAIL\t` under `-v`
  wrongly deflates outer `Passed` (e.g. `FAIL (1/2)` or `FAIL (0/2)`).
- Leaves are labeled `heavy`.

```go
import (
	"fmt"
	"os"
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

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")
var finalSummaryPlainRe = regexp.MustCompile(`^(PASS|FAIL) \(\d+/\d+\) in .+$`)

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

// createNestedFailOuterPassTree builds an outer 2-leaf tree:
//   - pass_leaf: always pass
//   - nested_fail_ok: runs nested doctest test on a sibling 1-fail child,
//     expects non-zero, prints nested stdout (contains FAIL () so the outer
//     go test -v stream leaks nested suite summaries.
// Returns the outer tree path (not the tmp parent).
func createNestedFailOuterPassTree(t *testing.T, bin string) string {
	t.Helper()
	if bin == "" {
		t.Fatal("bin is empty")
	}
	tmp := t.TempDir()
	child := filepath.Join(tmp, "child_fail")
	testtree.WritePassFailTree(t, child, 0, 1)

	outer := filepath.Join(tmp, "outer")
	// Embed absolute paths so the fixture leaf can shell out without env discovery.
	assertNested := fmt.Sprintf(`import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	cmd := exec.Command(%s, "test", "--no-color", %s)
	cmd.Env = append(os.Environ(),
		"GOCACHE="+t.TempDir(),
		"DOCTEST_LEAF_CACHE="+t.TempDir(),
		"DOCTEST_CACHE_HOME="+t.TempDir(),
	)
	out, runErr := cmd.CombinedOutput()
	// Leak nested summary into the outer go-test stream so text FAIL scanners
	// (the pre-fix verbose path) would see "FAIL (" / "FAIL\t".
	fmt.Print(string(out))
	if runErr == nil {
		t.Fatal("expected nested intentional fail to exit non-zero")
	}
	if !strings.Contains(string(out), "FAIL (") {
		t.Fatalf("expected nested FAIL ( summary in nested stdout, got:\n%%s", out)
	}
}
`, strconv.Quote(bin), strconv.Quote(child))

	testtree.WriteMinimalRunnableTree(t, outer, []testtree.LeafSpec{
		{
			Name:     "pass_leaf",
			Steps:    "always pass",
			Expected: "passes",
		},
		{
			Name:     "nested_fail_ok",
			Steps:    "run nested doctest that fails; expect non-zero; print nested stdout",
			Expected: "outer passes while nested FAIL ( appears in stream",
			AssertGo: assertNested,
		},
	})
	return outer
}

// createRealFailTree builds a 1-leaf tree that always Assert-fails.
func createRealFailTree(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WritePassFailTree(t, tmp, 0, 1)
	return tmp
}

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

// findResultSummary returns the last PASS/FAIL suite summary line on stdout.
// Nested intentional-fail leaves leak their own "FAIL (0/1)" into the stream;
// the outer harness summary is always the final one.
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

// combinedOutput joins stdout and stderr for stream assertions.
func combinedOutput(resp *Response) string {
	if resp == nil {
		return ""
	}
	return resp.Stdout + "\n" + resp.Stderr
}

// streamHasNestedFailSummary reports whether harness output contains a nested
// suite summary token FAIL ( — proof the intentional nested fail leaked.
func streamHasNestedFailSummary(combined string) bool {
	return strings.Contains(combined, "FAIL (")
}
```
