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
  redirects `os.Stdout` to a pipe, and calls `build.Test`.
- Backtick characters in embedded Go strings use `\x60` to avoid conflicting
  with the outer markdown code fence.

## Steps
1. Create a temp sub-tree with `PassCount` passing and `FailCount` failing leaves.
2. Optionally warm the go-test cache with a prior `build.Test` run (`WarmCache`).
3. Redirect `os.Stdout` to a pipe (non-TTY) and call `build.Test`.
4. Parse dots and summary from captured stdout into `Response`.

## Context
- Leaf names use `a_` / `z_` prefixes so pass packages sort before fail packages.
- Tests use `ColorAlways` / `ColorNever` for deterministic behavior without a
  real terminal. `ColorAuto` + pipe verifies auto-off behavior.
- ANSI detection uses the regex `\x1b\[[0-9;]*m`.

```go
import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
)

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")

var summaryFailField = regexp.MustCompile(`,\s*((?:\x1b\[[0-9;]*m)*)?(\d+ Fail)((?:\x1b\[[0-9;]*m)*)?`)

type Request struct {
	Color     core.ColorMode
	Count     int
	PassCount int
	FailCount int
	WarmCache bool
}

type Response struct {
	Output      string
	Dots        string
	Summary     string
	GenDir      string
	TestErr     error
}

func writeFile(t *testing.T, root, rel, content string) {
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

func writeLeaf(t *testing.T, root, name string, fail bool) {
	t.Helper()
	assertBody := "func Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"
	if fail {
		assertBody = "func Assert(t *testing.T, req *Request, resp *Response, err error) { t.Fatal(\"forced failure\") }\n"
	}
	writeFile(t, root, name+"/SETUP.md", ""+
		"## Steps\n"+
		"1. Leaf setup.\n\n"+
		"\x60\x60\x60go\n"+
		"import \"testing\"\n"+
		"func Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"+
		"\x60\x60\x60\n")
	writeFile(t, root, name+"/ASSERT.md", ""+
		"## Expected\n"+
		"- Leaf assertion.\n\n"+
		"\x60\x60\x60go\n"+
		assertBody+
		"\x60\x60\x60\n")
}

func createSubTree(t *testing.T, root string, passCount, failCount int) {
	t.Helper()

	writeFile(t, root, "SETUP.md", ""+
		"## Preconditions\n"+
		"- Minimal harness sub-tree.\n\n"+
		"## Steps\n"+
		"1. Run returns immediately.\n\n"+
		"\x60\x60\x60go\n"+
		"import \"testing\"\n"+
		"type Request struct{}\n"+
		"type Response struct{}\n"+
		"func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"+
		"\x60\x60\x60\n")

	for i := 0; i < passCount; i++ {
		name := "a_pass_" + itoa(i)
		writeLeaf(t, root, name, false)
	}
	for i := 0; i < failCount; i++ {
		name := "z_fail_" + itoa(i)
		writeLeaf(t, root, name, true)
	}

	if passCount == 0 && failCount == 0 {
		writeLeaf(t, root, "a_pass_0", false)
	}
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

func metricIsColored(summary, metric string) bool {
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

func metricIsPlain(summary, metric string) bool {
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

func Run(t *testing.T, req *Request) (*Response, error) {
	subRoot := t.TempDir()
	genDir := filepath.Join(t.TempDir(), "gendir")

	createSubTree(t, subRoot, req.PassCount, req.FailCount)

	opts := core.Options{
		GenDir:     genDir,
		RemoveTemp: false,
		Color:      req.Color,
		Count:      req.Count,
	}

	if req.WarmCache {
		if err := build.Test(subRoot, opts); err != nil {
			t.Fatalf("cache warmup run failed: %v", err)
		}
	}

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	testErr := build.Test(subRoot, opts)
	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	os.Stdout = oldStdout

	output := buf.String()
	dots, summary := splitDotsAndSummary(output)

	return &Response{
		Output:  output,
		Dots:    dots,
		Summary: summary,
		GenDir:  genDir,
		TestErr: testErr,
	}, nil
}
```