# Scenario

**Feature**: `build.Test` captures stdout and asserts ANSI coloring rules

```
# run generated go tests, emit progress
build.Test(dir, opts) -> go test -> dot per package -> summary line

# color mode
ColorAuto -> TTY check on stdout | ColorAlways -> force ANSI | ColorNever -> plain

# colored regions (non-verbose only)
fail dot -> red | pass dot -> plain | summary Pass/Fail/Cached -> green/red/gray
```

## Preconditions
- The `build` package is importable (`github.com/xhd2015/doctest/libdoc/build`).
- Each leaf configures `Request` fields; root `Run` builds a temp sub-tree,
  captures progress via `opts.Stdout` (`bytes.Buffer`, non-TTY), and calls
  `build.Test` — **never** swaps `os.Stdout` (parallel-safe under `t.Parallel()`).
- Backtick characters in embedded Go strings use `\x60` to avoid conflicting
  with the outer markdown code fence.

## Steps
1. Create a temp sub-tree with `PassCount` passing and `FailCount` failing leaves.
2. Optionally warm the go-test cache with a prior `build.Test` run (`WarmCache`).
3. Call `build.Test` with `opts.Stdout` set to a buffer (non-TTY → ColorAuto off).
4. Parse dots and summary from the buffer into `Response`.

## Context
- Leaf names use `a_` / `z_` prefixes so pass packages sort before fail packages.
- Tests use `ColorAlways` / `ColorNever` for deterministic behavior without a
  real terminal. `ColorAuto` + buffer verifies auto-off (non-file writer).
- ANSI detection uses the regex `\x1b\[[0-9;]*m`.
- `core.Options.Stdout` is required for parallel-safe harnesses; product already
  routes progress + color resolution through it (`build.Test`).

```go
import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")
var summaryFailField = regexp.MustCompile(`,\s*((?:\x1b\[[0-9;]*m)*)?(\d+ Fail)((?:\x1b\[[0-9;]*m)*)?`)
func writeFile(t *testing.T, root string, rel string, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(b[pos:])
}
func createSubTree(t *testing.T, root string, passCount int, failCount int) {
	t.Helper()
	testtree.WritePassFailTree(t, root, passCount, failCount)
}
func splitDotsAndSummary(output string) (string, string) {
	idx := strings.Index(output, "  (")
	if idx < 0 {
		return output, ""
	}
	return output[:idx], output[idx:]
}
func containsANSI(s string) bool {
	return ansiEscape.MatchString(s)
}
func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}
func metricIsColored(summary string, metric string) bool {
	idx := strings.Index(summary, metric)
	if idx < 0 {
		return false
	}
	start := idx - 20
	if start < 0 {
		start = 0
	}
	end := idx + len(metric) + 20
	if end > len(summary) {
		end = len(summary)
	}
	return containsANSI(summary[start:end])
}
func failFieldIsColored(summary string) bool {
	m := summaryFailField.FindString(summary)
	if m == "" {
		return false
	}
	return containsANSI(m)
}
func metricIsPlain(summary string, metric string) bool {
	idx := strings.Index(summary, metric)
	if idx < 0 {
		return false
	}
	start := idx - 4
	if start < 0 {
		start = 0
	}
	end := idx + len(metric) + 4
	if end > len(summary) {
		end = len(summary)
	}
	return !containsANSI(summary[start:end])
}
```
